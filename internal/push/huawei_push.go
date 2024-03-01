package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/cossim/go-hms-push/push/config"
	client "github.com/cossim/go-hms-push/push/core"
	"github.com/cossim/go-hms-push/push/model"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/notify"
	"log"
)

const (
	HIGH   = "high"
	NORMAL = "nornal"

	DefaultMaxRetry = 1

	DefaultAuthUrl = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"
	DefaultPushUrl = "https://push-api.cloud.huawei.com"
)

// HMSService 实现huawei推送
type HMSService struct {
	clients map[string]*client.HMSClient
}

func NewHMSService(cfg *config.Config) (*HMSService, error) {
	s := &HMSService{
		clients: make(map[string]*client.HMSClient),
	}

	var (
		AuthUrl = DefaultAuthUrl
		PushUrl = DefaultPushUrl
	)

	for _, v := range cfg.Huawei {
		if v.AuthUrl != "" {
			AuthUrl = v.AuthUrl
		}
		if v.PushUrl != "" {
			PushUrl = v.PushUrl
		}
		client, err := client.NewHttpClient(&c.Config{
			AppId:     v.AppID,
			AppSecret: v.AppSecret,
			AuthUrl:   AuthUrl,
			PushUrl:   PushUrl,
		})
		if err != nil {
			return nil, err
		}
		s.clients[v.AppID] = client
	}
	return s, nil
}

func (h *HMSService) Send(ctx context.Context, request interface{}) error {
	req, ok := request.(*notify.HMSPushNotification)
	if !ok {
		return errors.New("invalid request parameter")
	}

	var maxRetry = req.Retry
	if maxRetry <= 0 {
		maxRetry = DefaultMaxRetry // 设置一个默认的最大重试次数
	}

	if err := h.validation(req); err != nil {
		return err
	}

	client, ok := h.clients[req.AppID]
	if !ok {
		return errors.New("invalid appid or appid push is not enabled")
	}

	log.Printf("hms req %v", req)

	for retryCount := 0; retryCount < maxRetry; retryCount++ {
		notification, err := h.buildNotification(req)
		if err != nil {
			return err
		}
		marshal, err := json.Marshal(notification)
		if err != nil {
			return err
		}
		fmt.Println("marshal => ", string(marshal))
		res, err := client.SendMessage(ctx, notification)
		if err != nil {
			return err
		}

		if res.Code == "80000000" {
			log.Printf("Notification is sent successfully! Code: %s", res.Code)
			return nil
		}

		log.Printf("Huawei Send Notification is failed! Code: %s msg: %s", res.Code, res.Msg)
	}

	return fmt.Errorf("failed to send notification after %d attempts", maxRetry)
}

func (h *HMSService) MulticastSend(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (h *HMSService) Subscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (h *HMSService) Unsubscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (h *HMSService) SendToTopic(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (h *HMSService) SendToCondition(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (h *HMSService) CheckDevice(ctx context.Context, req interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (h *HMSService) Name() string {
	//TODO implement me
	panic("implement me")
}

func (h *HMSService) validation(req *notify.HMSPushNotification) error {
	if req.MessageRequest == nil {
		return errors.New("message request is empty")
	}
	return nil
}

// HTTP Connection Server Reference for HMS
// https://developer.huawei.com/consumer/en/doc/development/HMS-References/push-sendapi
func (h *HMSService) buildNotification(req *notify.HMSPushNotification) (*model.MessageRequest, error) {
	msgRequest := model.NewNotificationMsgRequest()

	msgRequest.Message.Android = model.GetDefaultAndroid()

	if len(req.Tokens) > 0 {
		msgRequest.Message.Token = req.Tokens
	}

	if len(req.Topic) > 0 {
		msgRequest.Message.Topic = req.Topic
	}

	if len(req.Condition) > 0 {
		msgRequest.Message.Condition = req.Condition
	}

	if req.Priority == HIGH {
		msgRequest.Message.Android.Urgency = "HIGH"
	}

	if len(req.Priority) > 0 {
		msgRequest.Message.Android.Urgency = req.Priority
	}

	if len(req.Category) > 0 {
		msgRequest.Message.Android.Category = req.Category
	}

	if len(req.TTL) > 0 {
		msgRequest.Message.Android.TTL = req.TTL
	}

	if len(req.BiTag) > 0 {
		msgRequest.Message.Android.BiTag = req.BiTag
	}

	msgRequest.Message.Android.FastAppTarget = req.FastAppTarget

	// Add data fields
	if len(req.Data) > 0 {
		msgRequest.Message.Data = req.Data
	}

	// Notification Message
	if req.MessageRequest.Message.Android.Notification != nil {
		msgRequest.Message.Android.Notification = req.MessageRequest.Message.Android.Notification

		if msgRequest.Message.Android.Notification.ClickAction == nil {
			msgRequest.Message.Android.Notification.ClickAction = model.GetDefaultClickAction()
		}
	}

	setDefaultAndroidNotification := func() {
		if msgRequest.Message.Android.Notification == nil {
			msgRequest.Message.Android.Notification = model.GetDefaultAndroidNotification()
		}
	}

	if len(req.Message) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Body = req.Message
	}

	if len(req.Title) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Title = req.Title
	}

	if len(req.Image) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Image = req.Image
	}

	if v, ok := req.Sound.(string); ok && len(v) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Sound = v
	} else if msgRequest.Message.Android.Notification != nil {
		msgRequest.Message.Android.Notification.DefaultSound = true
	}

	if req.Development {
		msgRequest.Message.Android.TargetUserType = 1
	}

	m, err := json.Marshal(msgRequest)
	if err != nil {
		log.Printf("Failed to marshal the default message! Error is " + err.Error())
		return nil, err
	}

	log.Printf("Default message is %s", string(m))
	return msgRequest, nil
}
