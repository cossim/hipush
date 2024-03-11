package push

import (
	"context"
	"errors"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/consts"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/status"
	"github.com/go-logr/logr"
	xp "github.com/yilee/xiaomi-push"
	"log"
)

var (
	MaxConcurrentXiaomiPushes = make(chan struct{}, 100)
)

// XiaomiPushService 小米推送 实现 PushService 接口
type XiaomiPushService struct {
	clients map[string]*xp.MiPush
	status  *status.StateStorage
	logger  logr.Logger
}

func NewXiaomiService(cfg *config.Config, logger logr.Logger) *XiaomiPushService {
	s := &XiaomiPushService{
		clients: map[string]*xp.MiPush{},
		status:  status.StatStorage,
		logger:  logger,
	}

	for _, v := range cfg.Xiaomi {
		if !v.Enabled || v.Enabled && v.AppSecret == "" {
			panic("push not enabled or misconfigured")
		}
		client := xp.NewClient(v.AppSecret, v.Package)
		s.clients[v.AppID] = client
	}

	return s
}

func (x *XiaomiPushService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.XiaomiPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	so := &SendOptions{}
	so.ApplyOptions(opt)

	if err := x.checkNotification(req); err != nil {
		return err
	}

	notification, err := x.buildNotification(req)
	if err != nil {
		return err
	}

	if so.DryRun {
		return nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return x.send(ctx, req.AppID, token, notification)
	}
	return RetrySend(ctx, send, req.Tokens, so.Retry, so.RetryInterval, 100)

	//for {
	//	newTokens, err := x.send(ctx, req.AppID, req.Tokens, notification)
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
	//fmt.Println("total count => ", x.status.GetXiaomiTotal())
	//fmt.Println("success count => ", x.status.GetXiaomiSuccess())
	//fmt.Println("failed count => ", x.status.GetXiaomiFailed())
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

func (x *XiaomiPushService) checkNotification(req *notify.XiaomiPushNotification) error {
	if len(req.Tokens) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.Title == "" {
		return errors.New("title cannot be empty")
	}

	if req.Content == "" {
		return errors.New("content cannot be empty")
	}

	if req.IsScheduled && req.ScheduledTime == 0 {
		return errors.New("scheduled time cannot be empty")
	}

	return nil
}

func (x *XiaomiPushService) buildNotification(req *notify.XiaomiPushNotification) (*xp.Message, error) {
	msg := xp.NewAndroidMessage(req.Title, req.Content).SetPayload("this is payload1")

	if req.NotifyType != 0 {
		msg.SetNotifyType(int32(req.NotifyType))
	}

	if req.TTL != 0 {
		msg.SetTimeToLive(req.TTL)
	}

	if req.IsScheduled && req.ScheduledTime != 0 {
		msg.SetTimeToSend(int64(req.ScheduledTime))
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
		log.Printf("xiaomi send success: %s", res.Reason)
		x.status.AddXiaomiSuccess(1)
		resp.Code = Success
		resp.Msg = res.Reason
	}

	return resp, err
}

func (x *XiaomiPushService) Name() string {
	return consts.PlatformXiaomi.String()
}
