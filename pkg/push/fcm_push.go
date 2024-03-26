package push

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-logr/logr"
	"google.golang.org/api/option"
	"log"
	"strings"
	"time"
)

var (
	_ push.PushService = &FCMService{}
)

// FCMService 谷歌安卓推送，实现了 PushService 接口
type FCMService struct {
	clients        map[string]*messaging.Client
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewFCMService(cfg *config.Config, logger logr.Logger) *FCMService {
	s := &FCMService{
		clients:        make(map[string]*messaging.Client),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
	}

	for _, v := range cfg.Android {
		if !v.Enabled {
			continue
		}
		if v.KeyPath == "" {
			panic("push not enabled or misconfigured")
		}

		opt := option.WithCredentialsFile(v.KeyPath)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			panic(err)
		}

		client, err := app.Messaging(context.Background())
		if err != nil {
			panic(err)
		}

		if v.AppID == "" {
			s.clients[v.AppName] = client
		} else {
			s.clients[v.AppID] = client
		}

		if v.AppName != "" {
			if v.AppID == "" {
				s.appNameToIDMap[v.AppName] = v.AppName
			} else {
				s.appNameToIDMap[v.AppName] = v.AppID
			}
		}
	}

	return s
}

func (f *FCMService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = f.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if err := f.checkNotification(req); err != nil {
		return nil, err
	}

	notification := f.buildAndroidNotification(req)

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return f.send(ctx, appid, token, notification)
	}

	retrySend, err := RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: retrySend.Data.(string)}, nil
}

func (f *FCMService) GetTasksStatus(ctx context.Context, appid string, taskID []string, obj push.TaskObjectList) error {
	return nil
}

func (f *FCMService) send(ctx context.Context, appid string, token string, notification *messaging.Message) (*Response, error) {
	client, ok := f.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	resp := &Response{Code: Fail}

	f.status.AddAndroidTotal(1)

	notification.Token = token
	res, err := client.Send(ctx, notification)
	if err != nil {
		log.Printf("fcm send error: %s", err)
		f.status.AddAndroidFailed(1)
		if res != "" {
			resp.Msg = res
		} else {
			resp.Msg = err.Error()
		}
	} else {
		log.Printf("fcm send success: %s", res)
		f.status.AddAndroidSuccess(1)
		resp.Code = Success
		resp.Msg = res
		resp.Data = res
	}

	return resp, err
}

// checkNotification for check request message
func (f *FCMService) checkNotification(req push.SendRequest) error {
	var msg string

	// ignore send topic mesaage from FCM
	if !f.IsTopic(req) && len(req.GetToken()) == 0 && req.GetTopic() == "" {
		msg = "the message must specify at least one registration ID"
		return errors.New(msg)
	}

	if len(req.GetToken()) == 0 {
		msg = "the token must not be empty"
		return errors.New(msg)
	}

	if len(req.GetToken()) > 1000 {
		msg = "the message may specify at most 1000 registration IDs"
		return errors.New(msg)
	}

	// ref: https://firebase.google.com/docs/cloud-messaging/http-server-ref
	ttlSeconds := int(req.GetTTL())
	if ttlSeconds < 0 || ttlSeconds > 2419200 {
		msg := "the message's TimeToLive field must be an integer " +
			"between 0 and 2419200 (4 weeks)"
		return errors.New(msg)
	}

	return nil
}

func (f *FCMService) IsTopic(req push.SendRequest) bool {
	return req.GetTopic() != "" && strings.HasPrefix(req.GetTopic(), "/topics/") && req.GetCondition() != ""
}

// buildAndroidNotification use for define Android notification.
// HTTP Connection Server Reference for Android
// https://firebase.google.com/docs/cloud-messaging/http-server-ref
func (f *FCMService) buildAndroidNotification(req push.SendRequest) *messaging.Message {
	notification := &messaging.Message{
		//Token:     req.Topic,
		Topic:     req.GetTopic(),
		Condition: req.GetCondition(),
		Android:   &messaging.AndroidConfig{},
	}

	if len(req.GetToken()) > 0 {
		//notification.Token = ""
		//notification.Token = req.Tokens[0]
	}

	if req.GetCategory() != "" {
		notification.Android.CollapseKey = req.GetCollapseID()
	}

	if req.GetPriority() == "high" || req.GetPriority() == "normal" {
		notification.Android.Priority = req.GetPriority()
	}

	// Add another field
	//if len(req.GetData()) > 0 {
	//	var data = make(map[string]string)
	//	for k, v := range req.GetData() {
	//		data[k] = fmt.Sprintf("%v", v)
	//	}
	//	notification.Data = data
	//}

	duration := time.Duration(req.GetTTL()) * time.Second
	notification.Android.TTL = &duration

	n := &messaging.Notification{}
	isNotificationSet := false

	if len(req.GetContent()) > 0 {
		isNotificationSet = true
		n.Body = req.GetContent()
	}

	if len(req.GetTitle()) > 0 {
		isNotificationSet = true
		n.Title = req.GetTitle()
	}

	if len(req.GetIcon()) > 0 {
		isNotificationSet = true
		n.ImageURL = req.GetIcon()
	}

	//if req.GetSound() != "" {
	//	isNotificationSet = true
	//	//n.Sound = req.Sound
	//	sound, ok := req.GetSound().(string)
	//	if ok {
	//		notification.Android.Notification.Sound = sound
	//	}
	//}

	if isNotificationSet {
		notification.Notification = n
	}

	// handle iOS apns in fcm
	//if len(req.Apns) > 0 {
	// Handle iOS APNS
	//}

	return notification
}

func (f *FCMService) Name() string {
	return consts.PlatformAndroid.String()
}
