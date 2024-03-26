package push

import (
	"context"
	"encoding/json"
	"errors"
	c "github.com/cossim/go-hms-push/push/config"
	hClient "github.com/cossim/go-hms-push/push/core"
	"github.com/cossim/go-hms-push/push/model"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-logr/logr"
	"log"
	"time"
)

const (
	HIGH   = "high"
	NORMAL = "nornal"

	DefaultAuthUrl = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"
	DefaultPushUrl = "https://push-api.cloud.huawei.com"
)

var (
	_ push.PushService = &HMSService{}
)

// HMSService 实现huawei推送，实现 PushService 接口
type HMSService struct {
	clients        map[string]*hClient.HMSClient
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewHMSService(cfg *config.Config, logger logr.Logger) *HMSService {
	s := &HMSService{
		clients:        make(map[string]*hClient.HMSClient),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
	}

	var (
		AuthUrl = DefaultAuthUrl
		PushUrl = DefaultPushUrl
	)

	for _, v := range cfg.Huawei {
		if !v.Enabled {
			continue
		}
		if v.AuthUrl != "" {
			AuthUrl = v.AuthUrl
		}
		if v.PushUrl != "" {
			PushUrl = v.PushUrl
		}
		if v.AppID == "" || v.AppSecret == "" {
			panic("invalid appid or appid push is not enabled")
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
		if v.AppName != "" {
			s.appNameToIDMap[v.AppName] = v.AppID
		}
	}

	return s
}

func (h *HMSService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = h.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if err := h.checkNotification(req); err != nil {
		return nil, err
	}

	notification, err := h.buildNotification(req, so)
	if err != nil {
		return nil, err
	}

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return h.send(ctx, appid, token, notification)
	}

	resp, err := RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	taskid, err := h.getTaskIDFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: taskid}, nil
}

// getTaskIDFromResponse 从 Response 结构体中获取 RequestId
func (h *HMSService) getTaskIDFromResponse(response *Response) (string, error) {
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		return "", err
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(marshal, &dataMap); err != nil {
		return "", err
	}
	requestID, ok := dataMap["requestId"].(string)
	if !ok {
		return "", errors.New("RequestId 字段不是 string 类型")
	}
	return requestID, nil
}

func (h *HMSService) GetTasksStatus(ctx context.Context, appid string, taskID []string, obj push.TaskObjectList) error {
	return nil
}

func (h *HMSService) send(ctx context.Context, appid string, token string, notification *model.MessageRequest) (*Response, error) {
	client, ok := h.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	h.status.AddHuaweiTotal(1)

	resp := &Response{}
	notification.Message.Token = []string{token}
	res, err := client.SendMessage(ctx, notification)
	if err != nil {
		log.Printf("huawei send error: %s", err)
		h.status.AddHuaweiFailed(1)
		resp.Code = Fail
		resp.Msg = res.Msg
	} else if res != nil && res.Code != "80000000" {
		log.Printf("huawei send error: %s", res.Msg)
		h.status.AddHonorFailed(1)
		err = errors.New(res.Msg)
		resp.Msg = res.Msg
	} else {
		log.Printf("huawei send success: %s", res)
		h.status.AddHuaweiSuccess(1)
		resp.Code = Success
		resp.Msg = res.Msg
		resp.Data = res
	}

	return resp, err
}

func (h *HMSService) checkNotification(req push.SendRequest) error {
	if len(req.GetToken()) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.GetTitle() == "" {
		return errors.New("title cannot be empty")
	}

	if req.GetContent() == "" {
		return errors.New("content cannot be empty")
	}

	return nil
}

// HTTP Connection Server Reference for HMS
// https://developer.huawei.com/consumer/en/doc/development/HMS-References/push-sendapi
func (h *HMSService) buildNotification(req push.SendRequest, so *push.SendOptions) (*model.MessageRequest, error) {
	msgRequest := model.NewNotificationMsgRequest()

	msgRequest.Message.Android = model.GetDefaultAndroid()

	if len(req.GetToken()) > 0 {
		msgRequest.Message.Token = req.GetToken()
	}

	if len(req.GetTopic()) > 0 {
		msgRequest.Message.Topic = req.GetTopic()
	}

	if len(req.GetCondition()) > 0 {
		msgRequest.Message.Condition = req.GetCondition()
	}

	if len(req.GetPriority()) > 0 {
		msgRequest.Message.Android.Urgency = req.GetPriority()
	}

	if len(req.GetCategory()) > 0 {
		msgRequest.Message.Android.Category = req.GetCategory()
	}

	if req.GetTTL() > 0 {
		duration := time.Duration(req.GetTTL()) * time.Second
		msgRequest.Message.Android.TTL = duration.String()
	}

	//if len(req.BiTag) > 0 {
	//	msgRequest.Message.Android.BiTag = req.BiTag
	//}

	//msgRequest.Message.Android.FastAppTarget = req.FastAppTarget

	// Add data fields
	//if len(req.Data) > 0 {
	//	msgRequest.Message.Data = req.Data
	//}

	// Notification Content
	//if req.MessageRequest.Message.Android.Notification != nil {
	//	msgRequest.Message.Android.Notification = req.MessageRequest.Message.Android.Notification
	//
	//	if msgRequest.Message.Android.Notification.ClickAction == nil {
	//		msgRequest.Message.Android.Notification.ClickAction = model.GetDefaultClickAction()
	//	}
	//}

	setDefaultAndroidNotification := func() {
		if msgRequest.Message.Android.Notification == nil {
			msgRequest.Message.Android.Notification = model.GetDefaultAndroidNotification()
		}
	}

	if len(req.GetContent()) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Body = req.GetContent()
	}

	if len(req.GetTitle()) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Title = req.GetTitle()
	}

	if len(req.GetIcon()) > 0 {
		setDefaultAndroidNotification()
		msgRequest.Message.Android.Notification.Image = req.GetIcon()
	}

	//if v, ok := req.Sound.(string); ok && len(v) > 0 {
	//	setDefaultAndroidNotification()
	//	msgRequest.Message.Android.Notification.Sound = v
	//} else if msgRequest.Message.Android.Notification != nil {
	//	msgRequest.Message.Android.Notification.DefaultSound = true
	//}
	msgRequest.Message.Android.Notification.DefaultSound = true

	if so.Development {
		msgRequest.Message.Android.TargetUserType = 1
	}

	return msgRequest, nil
}

func (h *HMSService) Name() string {
	return consts.PlatformHuawei.String()
}
