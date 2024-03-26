package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	xp "github.com/cossim/xiaomi-push"
	"github.com/go-logr/logr"
	"log"
	"strings"
)

var (
	_ push.PushService = &XiaomiPushService{}
)

// XiaomiPushService 小米推送 实现 PushService 接口
type XiaomiPushService struct {
	clients        map[string]*xp.MiPush
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewXiaomiService(cfg *config.Config, logger logr.Logger) *XiaomiPushService {
	s := &XiaomiPushService{
		clients:        make(map[string]*xp.MiPush),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
	}

	for _, v := range cfg.Xiaomi {
		if !v.Enabled || v.Enabled && v.AppSecret == "" {
			panic("push not enabled or misconfigured")
		}
		client := xp.NewClient(v.AppSecret, v.Package)
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

func (x *XiaomiPushService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = x.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if err := x.checkNotification(req); err != nil {
		return nil, err
	}

	notification, err := x.buildNotification(req)
	if err != nil {
		return nil, err
	}

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return x.send(ctx, appid, token, notification)
	}

	res, err := RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	taskid, err := x.getTaskIDFromResponse(res)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: taskid}, nil
}

// getTaskIDFromResponse 从 Response 结构体中获取 task_id 字段
func (x *XiaomiPushService) getTaskIDFromResponse(response *Response) (string, error) {
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		return "", err
	}
	var t map[string]interface{}
	if err := json.Unmarshal(marshal, &t); err != nil {
		return "", err
	}
	data, ok := t["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("data 字段不是 map[string]interface{} 类型")
	}
	taskid, ok := data["id"].(string)
	if !ok {
		return "", fmt.Errorf("id 字段不是 string 类型")
	}
	return taskid, nil
}

func (x *XiaomiPushService) GetTasksStatus(ctx context.Context, key string, taskID []string, list push.TaskObjectList) error {
	var appid string
	appid, ok := x.appNameToIDMap[key]
	if !ok {
		_, ok = x.clients[key]
		if !ok {
			return ErrInvalidAppID
		}
		appid = key
	}

	client, ok := x.clients[appid]
	if !ok {
		return ErrInvalidAppID
	}

	jobKey := strings.Join(taskID, ",")
	resp, err := client.GetMultiMessageStatusByMsgIDs(ctx, jobKey)
	if err != nil {
		x.logger.Error(err, "get tasks status failed")
		return err
	}

	if resp.Result.Code != 0 {
		x.logger.Error(errors.New(resp.Reason), "get tasks status failed")
		return errors.New(resp.Reason)
	}

	for _, result := range resp.Data.Data {
		obj := &push.VivoPushStats{}
		obj.SetTaskID(result.ID)
		obj.SetClick(int(result.Click))
		obj.SetDisplay(int(result.Delivered))
		obj.SetReceive(int(result.Resolved))
		obj.SetSend(int(result.Resolved))
		list.Add(obj)
	}

	return nil
}

func (x *XiaomiPushService) checkNotification(req push.SendRequest) error {
	if len(req.GetToken()) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.GetTitle() == "" {
		return errors.New("title cannot be empty")
	}

	if req.GetContent() == "" {
		return errors.New("content cannot be empty")
	}

	//if req.IsScheduled && req.ScheduledTime == 0 {
	//	return errors.New("scheduled time cannot be empty")
	//}

	return nil
}

func (x *XiaomiPushService) buildNotification(req push.SendRequest) (*xp.Message, error) {
	//msg := xp.NewAndroidMessage(req.Title, req.Content).SetPayload("this is payload1")
	msg := xp.NewAndroidMessage(req.GetTitle(), req.GetContent())
	if req.GetNotifyType() != 0 {
		msg.SetNotifyType(req.GetNotifyType())
	}

	if req.GetTTL() != 0 {
		msg.SetTimeToLive(req.GetTTL())
	}

	//if req.IsScheduled && req.ScheduledTime != 0 {
	//	msg.SetTimeToSend(int64(req.ScheduledTime))
	//}

	if req.GetForeground() {
		msg.Extra["notify_foreground"] = "1"
	} else {
		msg.Extra["notify_foreground"] = "0"
	}

	return msg, nil
}

func (x *XiaomiPushService) send(ctx context.Context, appID string, token string, message *xp.Message) (*Response, error) {
	client, ok := x.clients[appID]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	x.status.AddXiaomiTotal(1)

	resp := &Response{Code: Fail}
	res, err := client.Send(ctx, message, token)
	if err != nil {
		log.Printf("xiaomi send error: %s", err)
		x.status.AddXiaomiFailed(1)
		resp.Msg = err.Error()
	} else if res != nil && res.Code != 0 {
		log.Printf("xiaomi send error: %s", res.Reason)
		x.status.AddXiaomiFailed(1)
		err = errors.New(res.Reason)
		resp.Code = int(res.Code)
		resp.Msg = res.Reason
	} else {
		log.Printf("xiaomi send success: %v", res)
		x.status.AddXiaomiSuccess(1)
		resp.Code = Success
		resp.Msg = res.Reason
		resp.Data = res
	}
	return resp, err
}

func (x *XiaomiPushService) Name() string {
	return consts.PlatformXiaomi.String()
}
