package push

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/status"
	"github.com/go-logr/logr"
	"google.golang.org/api/option"
	"log"
	"strings"
	"sync"
)

var (
	// MaxConcurrentAndroidPushes pool to limit the number of concurrent iOS pushes
	MaxConcurrentAndroidPushes = make(chan struct{}, 100)
)

// FCMService 谷歌安卓推送，实现了 PushService 接口
type FCMService struct {
	clients map[string]*messaging.Client
	status  *status.StateStorage
	logger  logr.Logger
}

//func NewFCMService(cfg *config.Config) (*FCMService, error) {
//	s := &FCMService{
//		clients: make(map[string]*fcm.Client),
//	}
//	for _, v := range cfg.Android {
//		if v.AppKey == "" && v.Enabled {
//			return nil, errors.New("you should provide android.AppKey")
//		}
//		client, err := fcm.NewClient(v.AppKey)
//		if err != nil {
//			return nil, err
//		}
//		s.clients[v.AppID] = client
//	}
//
//	return s, nil
//}

func NewFCMService(cfg *config.Config, logger logr.Logger) (*FCMService, error) {
	s := &FCMService{
		clients: make(map[string]*messaging.Client),
		status:  status.StatStorage,
		logger:  logger,
	}

	for _, v := range cfg.Android {
		if !v.Enabled {
			continue
		}
		if v.Enabled && v.KeyPath == "" {
			return nil, errors.New("push not enabled or misconfigured")
		}

		opt := option.WithCredentialsFile(v.KeyPath)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			return nil, err
		}

		client, err := app.Messaging(context.Background())
		if err != nil {
			return nil, err
		}
		s.clients[v.AppID] = client
	}

	return s, nil
}

func (f *FCMService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.FCMPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	// 设置一个默认的最大重试次数
	var maxRetry = req.Retry
	var retryCount int
	var es []error
	if maxRetry <= 0 {
		maxRetry = DefaultMaxRetry
	}

	so := &SendOptions{}
	so.ApplyOptions(opt)

	if err := f.checkNotification(req); err != nil {
		return err
	}

	notification := f.buildAndroidNotification(req)

	if so.DryRun {
		return nil
	}

	for {
		newTokens, err := f.send(ctx, req.AppID, req.Tokens, notification)
		if err != nil {
			log.Printf("send error => %v", err)
			es = append(es, err)
		}
		// 如果有重试的 Token，并且未达到最大重试次数，则进行重试
		if len(newTokens) > 0 && retryCount < maxRetry {
			retryCount++
			req.Tokens = newTokens
		} else {
			break
		}
	}

	var errorMsgs []string
	for _, err := range es {
		errorMsgs = append(errorMsgs, err.Error())
	}
	if len(errorMsgs) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMsgs, ", "))
	}
	return nil
}

func (f *FCMService) send(ctx context.Context, appid string, tokens []string, notification *messaging.Message) ([]string, error) {
	var newTokens []string
	var wg sync.WaitGroup

	client, ok := f.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentAndroidPushes <- struct{}{}
		wg.Add(1)
		f.status.AddAndroidTotal(1)
		go func(notification *messaging.Message, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentAndroidPushes
				wg.Done()
			}()

			res, err := client.Send(ctx, notification)
			if err != nil {
				log.Printf("fcm send error: %s", err)
				newTokens = append(newTokens, token)
				f.status.AddAndroidFailed(1)
			} else {
				log.Printf("fcm send success: %s", res)
				f.status.AddAndroidSuccess(1)
			}
		}(notification, token)
	}
	wg.Wait()
	if len(es) > 0 {
		var errorStrings []string
		for _, err := range es {
			errorStrings = append(errorStrings, err.Error())
		}
		allErrorsString := strings.Join(errorStrings, ", ")
		return nil, errors.New(allErrorsString)
	}
	return newTokens, nil
}

// checkNotification for check request message
func (f *FCMService) checkNotification(req *notify.FCMPushNotification) error {
	var msg string

	// ignore send topic mesaage from FCM
	if !f.IsTopic(req) && len(req.Tokens) == 0 && req.Topic == "" {
		msg = "the message must specify at least one registration ID"
		return errors.New(msg)
	}

	if len(req.Tokens) == 0 {
		msg = "the token must not be empty"
		return errors.New(msg)
	}

	if len(req.Tokens) > 1000 {
		msg = "the message may specify at most 1000 registration IDs"
		return errors.New(msg)
	}

	// ref: https://firebase.google.com/docs/cloud-messaging/http-server-ref
	if req.TTL != nil && *req.TTL > uint(2419200) {
		msg = "the message's TimeToLive field must be an integer " +
			"between 0 and 2419200 (4 weeks)"
		return errors.New(msg)
	}

	return nil
}

func (f *FCMService) IsTopic(req *notify.FCMPushNotification) bool {
	return req.Topic != "" && strings.HasPrefix(req.Topic, "/topics/") && req.Condition != ""
}

// buildAndroidNotification use for define Android notification.
// HTTP Connection Server Reference for Android
// https://firebase.google.com/docs/cloud-messaging/http-server-ref
func (f *FCMService) buildAndroidNotification(req *notify.FCMPushNotification) *messaging.Message {
	notification := &messaging.Message{
		Token:     req.Topic,
		Condition: req.Condition,
	}

	if len(req.Tokens) > 0 {
		notification.Token = ""
		notification.Token = req.Tokens[0]
	}

	if req.Priority == HIGH || req.Priority == "normal" {
		notification.Android.Priority = req.Priority
	}

	// Add another field
	if len(req.Data) > 0 {
		//notification.Data = req.Data
	}

	n := &messaging.Notification{}
	isNotificationSet := false

	if len(req.Message) > 0 {
		isNotificationSet = true
		n.Body = req.Message
	}

	if len(req.Title) > 0 {
		isNotificationSet = true
		n.Title = req.Title
	}

	if len(req.Image) > 0 {
		isNotificationSet = true
		n.ImageURL = req.Image
	}

	if req.Sound != "" {
		isNotificationSet = true
		//n.Sound = req.Sound
		notification.Android.Notification.Sound = req.Sound
	}

	if isNotificationSet {
		notification.Notification = n
	}

	// handle iOS apns in fcm

	if len(req.Apns) > 0 {
		// Handle iOS APNS
	}

	return notification
}

func (f *FCMService) SendMulticast(ctx context.Context, req interface{}, opt ...MulticastOption) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) Subscribe(ctx context.Context, req interface{}, opt ...SubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) Unsubscribe(ctx context.Context, req interface{}, opt ...UnsubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) SendToTopic(ctx context.Context, req interface{}, opt ...TopicOption) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) CheckDevice(ctx context.Context, req interface{}, opt ...CheckDeviceOption) bool {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) Name() string {
	//TODO implement me
	panic("implement me")
}
