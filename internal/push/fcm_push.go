package push

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/notify"
	"google.golang.org/api/option"
	"log"
	"strings"
)

// FCMService 谷歌安卓推送
type FCMService struct {
	clients map[string]*messaging.Client
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

func NewFCMService(cfg *config.Config) (*FCMService, error) {
	s := &FCMService{
		clients: make(map[string]*messaging.Client),
	}
	fmt.Println("cfg.Android => ", cfg.Android)
	for _, v := range cfg.Android {
		fmt.Println(" v => ", v)
		if !v.Enabled || v.Enabled && v.KeyPath == "" {
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

func (f *FCMService) Send(ctx context.Context, request interface{}, opt SendOption) error {
	req, ok := request.(*notify.FCMPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	// 设置一个默认的最大重试次数
	var maxRetry = req.Retry
	if maxRetry <= 0 {
		maxRetry = DefaultMaxRetry
	}

	if err := f.validation(req); err != nil {
		return err
	}

	var retryCount int

	fmt.Println("req.AppID =?> ", req.AppID)

	client, ok := f.clients[req.AppID]
	if !ok {
		return errors.New("invalid appid or appid push is not enabled")
	}

Retry:
	notification := f.buildAndroidNotification(req)
	res, err := client.Send(ctx, notification)
	if err != nil {
		log.Printf("FCM server send message error: %s", err)
		return err
	}

	if !f.IsTopic(req) {
		log.Printf("FCM Success resp %v", res)
		//log.Printf(fmt.Sprintf("Android Success count: %d, Failure count: %d", res.Success, res.Failure))
	}

	// TODO 记录发送成功和失败的数据

	var newTokens []string
	// result from Send messages to specific devices
	//for k, result := range res.Results {
	//	to := req.Topic // 默认使用 Topic
	//	if k < len(req.Tokens) {
	//		to = req.Tokens[k]
	//	}
	//
	//	if result.Error != nil && !result.Unregistered() {
	//		newTokens = append(newTokens, to)
	//	}
	//}

	// result from Send messages to topics
	if f.IsTopic(req) {
		to := req.Topic
		if to == "" {
			to = req.Condition
		}
		log.Printf("Topic Content: %s", to)
	}

	// Device Group HTTP Response
	//if len(res.FailedRegistrationIDs) > 0 {
	//	newTokens = append(newTokens, res.FailedRegistrationIDs...)
	//}

	if len(newTokens) > 0 && retryCount < maxRetry {
		retryCount++
		req.Tokens = newTokens
		goto Retry
	}

	return nil
}

// validation for check request message
func (f *FCMService) validation(req *notify.FCMPushNotification) error {
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

func (f *FCMService) MulticastSend(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) Subscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) Unsubscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) SendToTopic(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) SendToCondition(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) CheckDevice(ctx context.Context, req interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (f *FCMService) Name() string {
	//TODO implement me
	panic("implement me")
}
