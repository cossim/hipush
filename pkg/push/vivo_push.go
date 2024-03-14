package push

import (
	"context"
	"errors"
	"fmt"
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
	clients map[string]*vp.VivoPush
	status  *status.StateStorage
	logger  logr.Logger

	scheduler gocron.Scheduler
}

func NewVivoService(cfg *config.Config, logger logr.Logger, scheduler gocron.Scheduler) *VivoService {
	s := &VivoService{
		clients:   map[string]*vp.VivoPush{},
		status:    status.StatStorage,
		logger:    logger,
		scheduler: scheduler,
	}

	for _, v := range cfg.Vivo {
		if !v.Enabled || v.Enabled && (v.AppID == "" || v.AppKey == "" || v.AppSecret == "") {
			panic("push not enabled or misconfigured")
		}
		client, err := vp.NewClient(v.AppID, v.AppKey, v.AppSecret)
		if err != nil {
			panic(err)
		}
		s.clients[v.AppID] = client
	}

	return s
}

func (v *VivoService) Send(ctx context.Context, appid string, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.VivoPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	so := &SendOptions{}
	so.ApplyOptions(opt)

	notification, err := v.buildNotification(req)
	if err != nil {
		return err
	}

	if so.DryRun {
		return nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return v.send(appid, token, notification)
	}

	if err := RetrySend(ctx, send, req.Tokens, so.Retry, so.RetryInterval, 100); err != nil {
		return err
	}

	_, err = v.scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()), gocron.NewTask(func(appid, taskID string) {
		fmt.Println("start vivo job")
		client, ok := v.clients[appid]
		if !ok {
			v.logger.Error(ErrInvalidAppID, ErrInvalidAppID.Error())
			return
		}
		resp, err := client.GetMessageStatusByJobKey(taskID)
		if err != nil {
			v.logger.Error(err, "vivo get status error")
			return
		}
		if resp.Result == 0 {
			for _, ss := range resp.Statistics {
				if ss.TaskId == taskID {
					v.status.SetVivoSend(int64(ss.Send))
					v.status.SetVivoReceive(int64(ss.Receive))
					v.status.SetVivoDisplay(int64(ss.Display))
					v.status.SetVivoClick(int64(ss.Click))
				}
			}
		}
	}, appid, req.TaskID))
	if err != nil {
		return err
	}

	return nil
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

func (v *VivoService) GetNotifyStatus(ctx context.Context, appid, notifyID string, obj NotifyObject) error {
	client, ok := v.clients[appid]
	if !ok {
		return ErrInvalidAppID
	}
	fmt.Println("notifyID => ", notifyID)
	key, err := client.GetMessageStatusByJobKey(notifyID)
	if err != nil {
		fmt.Println("VivoService GetNotifyStatus err => ", err)
		return err
	}
	fmt.Println("VivoService GetNotifyStatus => ", key)
	return err
}

func (v *VivoService) Name() string {
	return consts.PlatformVivo.String()
}
