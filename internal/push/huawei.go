package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cossim/hipush/internal/notify"
	c "github.com/msalihkarakasli/go-hms-push/push/config"
	client "github.com/msalihkarakasli/go-hms-push/push/core"
	"github.com/msalihkarakasli/go-hms-push/push/model"
	"log"
)

const (
	HIGH   = "high"
	NORMAL = "nornal"

	DefaultMaxRetry = 5
)

// HMSService 实现huawei推送
type HMSService struct {
	client *client.HMSClient
}

func (h *HMSService) Send(ctx context.Context, request interface{}) error {
	req, ok := request.(notify.HMSPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	var maxRetry = req.Retry
	if maxRetry <= 0 {
		maxRetry = DefaultMaxRetry // 设置一个默认的最大重试次数
	}

	if err := h.validation(&req); err != nil {
		return err
	}

	for retryCount := 0; retryCount < maxRetry; retryCount++ {
		res, err := h.client.SendMessage(ctx, req.MessageRequest)
		if err != nil {
			return err
		}

		if res.Code == "80000000" {
			log.Printf("Notification is sent successfully! Code: %s", res.Code)
			return nil
		}

		log.Printf("Huawei Send Notification is failed! Code: %s", res.Code)
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
		return errors.New("invalid request")
	}
	return nil
}

// HTTP Connection Server Reference for HMS
// https://developer.huawei.com/consumer/en/doc/development/HMS-References/push-sendapi
func (h *HMSService) getNotification(req *notify.HMSPushNotification) (*model.MessageRequest, error) {
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

	m, err := json.Marshal(msgRequest)
	if err != nil {
		log.Printf("Failed to marshal the default message! Error is " + err.Error())
		return nil, err
	}

	log.Printf("Default message is %s", string(m))
	return msgRequest, nil
}

func NewHMSService(appID, appSecret string) *HMSService {
	conf := &c.Config{
		AppId:     appID,
		AppSecret: appSecret,
		AuthUrl:   "https://oauth-login.cloud.huawei.com/oauth2/v3/token",
		PushUrl:   "https://push-api.cloud.huawei.com",
	}
	client, err := client.NewHttpClient(conf)
	if err != nil {
		panic(err)
	}
	return &HMSService{client: client}
}
