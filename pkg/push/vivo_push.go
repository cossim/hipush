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
	vp "github.com/cossim/vivo-push"
	"github.com/go-logr/logr"
	"log"
	"net/url"
	"strings"
)

var (
	_ push.PushService = &VivoService{}
)

// VivoService 实现vivo推送，实现 PushService 接口
type VivoService struct {
	clients        map[string]*vp.VivoPush
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewVivoService(cfg *config.Config, logger logr.Logger) *VivoService {
	s := &VivoService{
		clients:        make(map[string]*vp.VivoPush),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
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

func (v *VivoService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = v.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	notification, err := v.buildNotification(req, so)
	if err != nil {
		return nil, err
	}

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return v.send(appid, token, notification)
	}

	resp, err := RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	taskid, err := v.getTaskIDFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: taskid}, nil
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

func (v *VivoService) buildNotification(req push.SendRequest, so *push.SendOptions) (*vp.Message, error) {
	// 检查 tokens 是否为空
	if len(req.GetToken()) == 0 {
		return nil, errors.New("tokens cannot be empty")
	}

	if req.GetTitle() == "" {
		return nil, errors.New("title cannot be empty")
	}

	if req.GetContent() == "" {
		return nil, errors.New("content cannot be empty")
	}

	// 检查 ClickAction 是否为空，为空则使用默认值
	clickAction := req.GetClickAction()
	if clickAction == nil {
		clickAction.Action = 1
	}

	if clickAction.Action == 0 {
		// 设置默认的 ClickAction
		clickAction.Action = 1
	}

	// 检查 URL 是否为合法 URL
	if clickAction.Action == 2 {
		_, err := url.Parse(clickAction.Url)
		if err != nil {
			return nil, err
		}
	}

	var notifyType = req.GetNotifyType()
	// 检查 NotifyType 是否为有效值
	if notifyType != 0 && notifyType < 1 || notifyType > 4 {
		return nil, errors.New("invalid notify type")
	}

	if req.GetNotifyType() == 0 {
		notifyType = 2
	}

	var pushMode int
	if so.Development {
		pushMode = 1
	}

	var ttl int64
	if req.GetTTL() == 0 {
		ttl = 60
	}

	data := make(map[string]string)
	for key, value := range req.GetCustomData() {
		if strValue, ok := value.(string); ok {
			data[key] = strValue
		} else {
			data[key] = fmt.Sprintf("%v", value)
		}
	}

	message := &vp.Message{
		RegId:           strings.Join(req.GetToken(), ","),
		NotifyType:      int(notifyType),
		Title:           req.GetTitle(),
		Content:         req.GetContent(),
		TimeToLive:      ttl,
		SkipType:        1,
		SkipContent:     req.GetClickAction().Url,
		NetworkType:     -1,
		ClientCustomMap: data,
		//Extra:           req.Data.ExtraMap(),
		//RequestId:      req.RequestId,
		//NotifyID:       req.NotifyID,
		Category:       req.GetCategory(),
		PushMode:       pushMode, // 默认为正式推送
		ForegroundShow: req.GetForeground(),
	}
	return message, nil
}

func (v *VivoService) GetTasksStatus(ctx context.Context, key string, tasks []string, list push.TaskObjectList) error {
	var appid string
	appid, ok := v.appNameToIDMap[key]
	if !ok {
		_, ok = v.clients[key]
		if !ok {
			return ErrInvalidAppID
		}
		appid = key
	}

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
