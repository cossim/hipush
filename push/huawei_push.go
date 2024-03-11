package push

import (
	"context"
	"encoding/json"
	"errors"
	c "github.com/cossim/go-hms-push/push/config"
	hClient "github.com/cossim/go-hms-push/push/core"
	"github.com/cossim/go-hms-push/push/model"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/consts"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/status"
	"github.com/go-logr/logr"
	"log"
)

const (
	HIGH   = "high"
	NORMAL = "nornal"

	DefaultMaxRetry = 1

	DefaultAuthUrl = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"
	DefaultPushUrl = "https://push-api.cloud.huawei.com"
)

var (
	// MaxConcurrentHuaweiPushes pool to limit the number of concurrent iOS pushes
	MaxConcurrentHuaweiPushes = make(chan struct{}, 100)
)

// HMSService 实现huawei推送，实现 PushService 接口
type HMSService struct {
	clients map[string]*hClient.HMSClient
	status  *status.StateStorage
	logger  logr.Logger
}

func NewHMSService(cfg *config.Config, logger logr.Logger) *HMSService {
	s := &HMSService{
		clients: make(map[string]*hClient.HMSClient),
		status:  status.StatStorage,
		logger:  logger,
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
		client, err := hClient.NewHttpClient(&c.Config{
			AppId:     v.AppID,
			AppSecret: v.AppSecret,
			AuthUrl:   AuthUrl,
			PushUrl:   PushUrl,
		})
		if err != nil {
			panic(err)
		}
		s.clients[v.AppID] = client
	}

	return s
}

func (h *HMSService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.HMSPushNotification)
	if !ok {
		return errors.New("invalid request parameter")
	}

	so := &SendOptions{}
	so.ApplyOptions(opt)

	if err := h.checkNotification(req); err != nil {
		return err
	}

	notification, err := h.buildNotification(req)
	if err != nil {
		return err
	}

	if so.DryRun {
		return nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return h.send(ctx, req.AppID, token, notification)
	}
	return RetrySend(ctx, send, req.Tokens, so.Retry, so.RetryInterval, 100)

	//for {
	//	newTokens, err := h.send(ctx, req.AppID, req.Tokens, notification)
	//	if err != nil {
	//		log.Printf("send error => %v", err)
	//		es = append(es, err)
	//	}
	//	// 如果有重试的 Token，并且未达到最大重试次数，则进行重试
	//	if len(newTokens) > 0 && retryCount < maxRetry {
	//		retryCount++
	//		req.Tokens = newTokens
	//	} else {
	//		break
	//	}
	//}
	//
	//var errorMsgs []string
	//for _, err := range es {
	//	errorMsgs = append(errorMsgs, err.Error())
	//}
	//if len(errorMsgs) > 0 {
	//	return fmt.Errorf("%s", strings.Join(errorMsgs, ", "))
	//}
	//
	//return nil
}

func (h *HMSService) send(ctx context.Context, appid string, token string, notification *model.MessageRequest) (*Response, error) {
	client, ok := h.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	resp := &Response{}
	notification.Message.Token = []string{token}
	res, err := client.SendMessage(ctx, notification)
	if err != nil {
		log.Printf("hcm send error: %s", err)
		h.status.AddHuaweiFailed(1)
		resp.Code = Fail
		resp.Msg = res.Msg
	} else if res != nil && res.Code != "80000000" {
		log.Printf("honor send error: %s", res.Msg)
		h.status.AddHonorFailed(1)
		err = errors.New(res.Msg)
		resp.Msg = res.Msg
	} else {
		log.Printf("hcm send success: %s", res)
		h.status.AddHuaweiSuccess(1)
		resp.Code = Success
		resp.Msg = res.Msg
	}

	return resp, err
}

func (h *HMSService) checkNotification(req *notify.HMSPushNotification) error {
	if len(req.Tokens) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.Title == "" {
		return errors.New("title cannot be empty")
	}

	if req.Content == "" {
		return errors.New("content cannot be empty")
	}

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

	// Notification Content
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

	if len(req.Content) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Body = req.Content
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

func (h *HMSService) Name() string {
	return consts.PlatformHuawei.String()
}
