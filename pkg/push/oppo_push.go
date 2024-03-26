package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	op "github.com/316014408/oppo-push"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-logr/logr"
	"github.com/golang/protobuf/jsonpb"
	"log"
)

var (
	_ push.PushService = &OppoService{}
)

// OppoService 实现oppo推送，实现 PushService 接口
type OppoService struct {
	clients        map[string]*op.OppoPush
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewOppoService(cfg *config.Config, logger logr.Logger) *OppoService {
	s := &OppoService{
		clients:        map[string]*op.OppoPush{},
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
	}

	for _, v := range cfg.Oppo {
		if !v.Enabled {
			continue
		}
		if v.AppKey == "" || v.AppSecret == "" {
			panic("push not enabled or misconfigured")
		}
		client := op.NewClient(v.AppKey, v.AppSecret)
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

func (o *OppoService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = o.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if err := o.checkNotification(req); err != nil {
		return nil, err
	}

	notification, err := o.buildNotification(req)
	if err != nil {
		return nil, err
	}

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return o.send(appid, token, notification)
	}

	resp, err := RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	taskid, err := o.getTaskIDFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: taskid}, nil
}

// getTaskIDFromResponse 从 Response 结构体中获取 RequestId
func (o *OppoService) getTaskIDFromResponse(response *Response) (string, error) {
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
	taskid, ok := data["messageId"].(string)
	if !ok {
		return "", fmt.Errorf("id 字段不是 string 类型")
	}
	return taskid, nil
}

func (o *OppoService) GetTasksStatus(ctx context.Context, appid string, taskID []string, obj push.TaskObjectList) error {
	return nil
}

func (o *OppoService) send(appID string, token string, notification *op.Message) (*Response, error) {
	client, ok := o.clients[appID]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	o.status.AddOppoTotal(1)

	resp := &Response{Code: Fail}
	notification.SetTargetValue(token)
	res, err := client.Unicast(notification)
	if err != nil {
		log.Printf("oppo send error: %s", err)
		o.status.AddOppoFailed(1)
		resp.Msg = err.Error()
	} else if res != nil && res.Code != 0 {
		log.Printf("oppo send error: %s", res.Message)
		o.status.AddOppoFailed(1)
		err = errors.New(res.Message)
		resp.Code = res.Code
		resp.Msg = res.Message
	} else {
		log.Printf("oppo send success: %v", res)
		o.status.AddOppoSuccess(1)
		resp.Code = Success
		resp.Msg = res.Message
		resp.Data = res
	}

	return resp, err
}

func (o *OppoService) checkNotification(req push.SendRequest) error {
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

func (o *OppoService) buildNotification(req push.SendRequest) (*op.Message, error) {
	m := op.NewMessage(req.GetTitle(), req.GetContent()).
		//SetSubTitle(req.Subtitle).
		SetTargetType(2)
	if req.GetClickAction() != nil {
		if req.GetClickAction().Action == 1 || req.GetClickAction().Action == 4 {
			m.SetClickActionActivity(req.GetClickAction().Activity)
		}
		if req.GetClickAction().Action == 2 {
			m.SetClickActionUrl(req.GetClickAction().Url)
		}
		m.SetClickActionType(int(req.GetClickAction().Action))
		marshaler := &jsonpb.Marshaler{}
		jsonString, err := marshaler.MarshalToString(req.GetClickAction().Parameters)
		if err != nil {
			return nil, err
		}
		m.SetActionParameters(jsonString)
	}

	return m, nil
}

func (o *OppoService) Name() string {
	return consts.PlatformOppo.String()
}
