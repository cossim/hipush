package push

import (
	"context"
	"encoding/json"
	"errors"
	mzp "github.com/cossim/go-meizu-push-sdk"
	v1 "github.com/cossim/hipush/api/pb/v1"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-logr/logr"
	"log"
)

var (
	_ push.PushService = &MeizuService{}
)

// MeizuService 实现魅族推送，实现 PushService 接口
type MeizuService struct {
	clients        map[string]func(token, message string) mzp.PushResponse
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewMeizuService(cfg *config.Config, logger logr.Logger) *MeizuService {
	s := &MeizuService{
		clients:        make(map[string]func(token, message string) mzp.PushResponse),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
	}

	for _, v := range cfg.Meizu {
		if !v.Enabled {
			continue
		}
		if v.AppID == "" || v.AppKey == "" {
			panic("push not enabled or misconfigured")
		}
		s.clients[v.AppID] = func(token, message string) mzp.PushResponse {
			appid := v.AppID
			appkey := v.AppKey
			return mzp.PushNotificationMessageByPushId(appid, token, message, appkey)
		}
		if v.AppName != "" {
			s.appNameToIDMap[v.AppName] = v.AppID
		}
	}

	return s
}

func (m *MeizuService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = m.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if err := m.checkNotification(req); err != nil {
		return nil, err
	}

	notification, err := m.buildNotification(req)
	if err != nil {
		return nil, err
	}

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return m.send(appid, token, notification)
	}

	_, err = RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{}, nil
}

func (m *MeizuService) GetTasksStatus(ctx context.Context, appid string, taskID []string, obj push.TaskObjectList) error {
	return nil
}

func (m *MeizuService) send(appid string, token string, message string) (*Response, error) {
	pushFunc, ok := m.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	m.status.AddMeizuTotal(1)

	var err error
	resp := &Response{}
	res := pushFunc(token, message)
	if res.GetCode() != Success {
		log.Printf("meizu send error: %s", res.GetMessage())
		m.status.AddMeizuFailed(1)
		err = errors.New(res.GetMessage())
		resp.Code = Fail
		resp.Msg = res.GetMessage()
	} else {
		log.Printf("meizu send success code: %v msg: %s", res.GetCode(), res.GetMessage())
		m.status.AddMeizuSuccess(1)
		resp.Code = Success
		resp.Msg = res.GetMessage()
		resp.Data = res
	}

	return resp, err
}

func (m *MeizuService) checkNotification(req push.SendRequest) error {
	if len(req.GetToken()) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.GetTitle() == "" {
		return errors.New("title cannot be empty")
	}

	if req.GetContent() == "" {
		return errors.New("message cannot be empty")
	}

	//if req.IsScheduled && (req.ScheduledStartTime == "" || req.ScheduledEndTime == "") {
	//	return errors.New("scheduled time cannot be empty")
	//}

	return nil
}

func (m *MeizuService) buildNotification(req push.SendRequest) (string, error) {
	msg := mzp.BuildNotificationMessage()
	msg.NoticeBarInfo.Title = req.GetTitle()
	msg.NoticeBarInfo.Content = req.GetContent()
	msg.ClickTypeInfo = mzp.ClickTypeInfo{
		ClickType:  int(req.GetClickAction().Action),
		Url:        req.GetClickAction().Url,
		Parameters: v1.StructPBToMap(req.GetClickAction().Parameters),
		Activity:   req.GetClickAction().Activity,
	}

	//offLine := 0
	//if req.OffLine {
	//	offLine = 1
	//}
	msg.PushTimeInfo = mzp.PushTimeInfo{
		//OffLine:   offLine,
		ValidTime: int(req.GetTTL()),
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(message), nil
}

func (m *MeizuService) Name() string {
	return consts.PlatformMeizu.String()
}
