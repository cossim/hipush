package push

import (
	"context"
	"errors"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/consts"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/status"
	vp "github.com/cossim/vivo-push"
	"github.com/go-logr/logr"
	"log"
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
}

func NewVivoService(cfg *config.Config, logger logr.Logger) *VivoService {
	s := &VivoService{
		clients: map[string]*vp.VivoPush{},
		status:  status.StatStorage,
		logger:  logger,
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

func (v *VivoService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
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
		return v.send(req.AppID, token, notification)
	}
	return RetrySend(ctx, send, req.Tokens, so.Retry, so.RetryInterval, 100)

	//for {
	//	newTokens, err := v.send(req)
	//	if err != nil {
	//		log.Printf("send error => %v", err)
	//		es = append(es, err)
	//	}
	//	// 如果有重试的 Token，并且未达到最大重试次数，则进行重试
	//	fmt.Println("retryCount => ", retryCount)
	//	fmt.Println("maxRetry => ", maxRetry)
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
	//return nil
}

func (v *VivoService) send(appid string, token string, notification *vp.Message) (*Response, error) {
	client, ok := v.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	resp := &Response{
		Code: Fail,
	}
	notification.RegId = token
	res, err := client.Send(notification, token)
	if err != nil {
		log.Printf("oppo send error: %s", err)
		v.status.AddOppoFailed(1)
		resp.Msg = err.Error()
	} else if res != nil && res.Result != 0 {
		log.Printf("oppo send error: %s", res.Desc)
		v.status.AddOppoFailed(1)
		err = errors.New(res.Desc)
		resp.Code = res.Result
		resp.Msg = res.Desc
	} else {
		log.Printf("oppo send success: %s", res.Desc)
		v.status.AddOppoSuccess(1)
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
		Action:  1,
		Content: "",
	}

	// 检查 ClickAction 是否为空，为空则使用默认值
	clickAction := req.ClickAction
	if clickAction == nil {
		clickAction = defaultClickAction
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
		SkipContent:     clickAction.Content,
		NetworkType:     -1,
		ClientCustomMap: req.Data,
		//Extra:           req.Data.ExtraMap(),
		RequestId:      req.RequestId,
		PushMode:       pushMode, // 默认为正式推送
		ForegroundShow: req.Foreground,
	}
	return message, nil
}

func (v *VivoService) Name() string {
	return consts.PlatformVivo.String()
}
