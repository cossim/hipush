package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/notify"
	"github.com/cossim/hipush/pkg/status"
	vp "github.com/cossim/vivo-push"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-logr/logr"
	"log"
	"net/url"
	"strings"
)

var (
	// MaxConcurrentVivoPushes pool to limit the number of concurrent iOS pushes
	MaxConcurrentVivoPushes = make(chan struct{}, 100)
)

// VivoService 实现vivo推送，实现 PushService 接口
type VivoService struct {
	clients        map[string]*vp.VivoPush
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger

	scheduler gocron.Scheduler
}

func NewVivoService(cfg *config.Config, logger logr.Logger, scheduler gocron.Scheduler) *VivoService {
	s := &VivoService{
		clients:        make(map[string]*vp.VivoPush),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
		scheduler:      scheduler,
	}

	for _, v := range cfg.Vivo {
		if !v.Enabled {
			continue
		}
		if v.AppID == "" || v.AppKey == "" || v.AppSecret == "" {
			panic("push not enabled or misconfigured")
		}
		client, err := vp.NewClient(v.AppID, v.AppKey, v.AppSecret)
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

func (v *VivoService) Send(ctx context.Context, request interface{}, opt ...push.SendOption) (*push.SendResponse, error) {
	req, ok := request.(*notify.VivoPushNotification)
	if !ok {
		return nil, errors.New("invalid request")
	}

	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	notification, err := v.buildNotification(req)
	if err != nil {
		return nil, err
	}

	var appid string
	if req.AppID != "" {
		appid = req.AppID
	} else if req.AppName != "" {
		appid, ok = v.appNameToIDMap[req.AppName]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return v.send(appid, token, notification)
	}

	resp, err := RetrySend(ctx, send, req.Tokens, so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	taskid, err := v.getTaskIDFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: taskid}, nil
	//
	//if err := RetrySend(ctx, send, req.Tokens, so.Retry, so.RetryInterval, 100); err != nil {
	//	return nil,err
	//}
	//
	//_, err = v.scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()), gocron.NewTask(func(appid, taskID string) {
	//	fmt.Println("start vivo job")
	//	client, ok := v.clients[appid]
	//	if !ok {
	//		v.logger.Error(ErrInvalidAppID, ErrInvalidAppID.Error())
	//		return
	//	}
	//	resp, err := client.GetMessageStatusByJobKey(taskID)
	//	if err != nil {
	//		v.logger.Error(err, "vivo get status error")
	//		return
	//	}
	//	if resp.Result == 0 {
	//		for _, ss := range resp.Statistics {
	//			if ss.TaskId == taskID {
	//				v.status.SetVivoSend(int64(ss.Send))
	//				v.status.SetVivoReceive(int64(ss.Receive))
	//				v.status.SetVivoDisplay(int64(ss.Display))
	//				v.status.SetVivoClick(int64(ss.Click))
	//			}
	//		}
	//	}
	//}, appid, req.TaskID))
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}

// getTaskIDFromResponse 从 Response 结构体中获取 task_id 字段
func (v *VivoService) getTaskIDFromResponse(response *Response) (string, error) {
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		return "", err
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(marshal, &dataMap); err != nil {
		return "", err
	}
	requestID, ok := dataMap["taskId"].(string)
	if !ok {
		return "", errors.New("RequestId 字段不是 string 类型")
	}
	return requestID, nil
}

func (v *VivoService) send(appid string, token string, notification *vp.Message) (*Response, error) {
	client, ok := v.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	v.status.AddVivoTotal(1)

	fmt.Println("notification Category => ", notification.Category)

	resp := &Response{Code: Fail}
	notification.RegId = token
	res, err := client.Send(notification, token)
	if err != nil {
		log.Printf("vivo send error: %s", err)
		v.status.AddVivoFailed(1)
		resp.Msg = err.Error()
	} else if res != nil && res.Result != 0 {
		log.Printf("vivo send error: %s", res.Desc)
		v.status.AddVivoFailed(1)
		err = errors.New(res.Desc)
		resp.Code = res.Result
		resp.Msg = res.Desc
	} else {
		log.Printf("vivo send success taskId: %v", res.TaskId)
		v.status.AddVivoSuccess(1)
		resp.Code = Success
		resp.Msg = res.Desc
		resp.Data = res
	}

	return resp, err
}

func (v *VivoService) buildNotification(req *notify.VivoPushNotification) (*vp.Message, error) {
	// 检查 tokens 是否为空
	if len(req.Tokens) == 0 {
		return nil, errors.New("tokens cannot be empty")
	}

	if req.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	if req.Message == "" {
		return nil, errors.New("message cannot be empty")
	}

	// 设置默认的 ClickAction
	defaultClickAction := &notify.VivoClickAction{
		Action: 1,
	}

	// 检查 ClickAction 是否为空，为空则使用默认值
	clickAction := req.ClickAction
	if clickAction == nil {
		clickAction = defaultClickAction
	}

	if clickAction.Action == 0 {
		clickAction.Action = 1
	}

	// 检查 URL 是否为合法 URL
	if clickAction.Action == 2 {
		_, err := url.Parse(clickAction.Url)
		if err != nil {
			return nil, err
		}
	}

	// 检查 NotifyType 是否为有效值
	if req.NotifyType != 0 && req.NotifyType < 1 || req.NotifyType > 4 {
		return nil, errors.New("invalid notify type")
	}

	if req.NotifyType == 0 {
		req.NotifyType = 2
	}

	var pushMode int
	if req.Development {
		pushMode = 1
	}

	if req.TTL == 0 {
		req.TTL = 60
	}

	message := &vp.Message{
		RegId:           strings.Join(req.Tokens, ","),
		NotifyType:      req.NotifyType,
		Title:           req.Title,
		Content:         req.Message,
		TimeToLive:      int64(req.TTL),
		SkipType:        clickAction.Action,
		SkipContent:     clickAction.Url,
		NetworkType:     -1,
		ClientCustomMap: req.Data,
		//Extra:           req.Data.ExtraMap(),
		RequestId:      req.RequestId,
		NotifyID:       req.NotifyID,
		Category:       req.Category,
		PushMode:       pushMode, // 默认为正式推送
		ForegroundShow: req.Foreground,
	}
	return message, nil
}

func (v *VivoService) GetTasksStatus(ctx context.Context, appid string, tasks []string, list push.TaskObjectList) error {
	client, ok := v.clients[appid]
	if !ok {
		return ErrInvalidAppID
	}

	jobKey := strings.Join(tasks, ",")
	resp, err := client.GetMessageStatusByJobKey(jobKey)
	if err != nil {
		return err
	}
	if resp.Result != 0 {
		return errors.New(resp.Desc)
	}

	for _, result := range resp.Statistics {
		obj := &push.VivoPushStats{}
		obj.SetTaskID(result.TaskId)
		obj.SetInvalidDevice(result.TargetInvalid)
		obj.SetClick(result.Click)
		obj.SetDisplay(result.Display)
		obj.SetReceive(result.Receive)
		obj.SetSend(result.Send)
		list.Add(obj)
	}

	return err
}

func (v *VivoService) Name() string {
	return consts.PlatformVivo.String()
}
