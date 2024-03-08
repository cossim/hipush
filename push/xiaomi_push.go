package push

import (
	"context"
	"errors"
	"fmt"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/status"
	"github.com/go-logr/logr"
	xp "github.com/yilee/xiaomi-push"
	"log"
	"strings"
	"sync"
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

func NewXiaomiService(cfg *config.Config, logger logr.Logger) (*XiaomiPushService, error) {
	s := &XiaomiPushService{
		clients: map[string]*xp.MiPush{},
		status:  status.StatStorage,
		logger:  logger,
	}

	for _, v := range cfg.Xiaomi {
		if !v.Enabled || v.Enabled && v.AppSecret == "" {
			return nil, errors.New("push not enabled or misconfigured")
		}
		client := xp.NewClient(v.AppSecret, v.Package)
		s.clients[v.AppID] = client
	}

	return s, nil
}

func (x *XiaomiPushService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.XiaomiPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	so := &SendOptions{}
	so.ApplyOptions(opt)

	var (
		retry      = so.Retry
		maxRetry   = retry
		retryCount = 0
		es         []error
	)

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

	for {
		newTokens, err := x.send(ctx, req.AppID, req.Tokens, notification)
		if err != nil {
			log.Printf("send error => %v", err)
			es = append(es, err)
		}
		// 如果有重试的 Token，并且未达到最大重试次数，则进行重试
		if len(newTokens) > 0 && retryCount < maxRetry {
			retryCount++
			req.Tokens = newTokens
		} else {
			break
		}
	}

	fmt.Println("total count => ", x.status.GetXiaomiTotal())
	fmt.Println("success count => ", x.status.GetXiaomiSuccess())
	fmt.Println("failed count => ", x.status.GetXiaomiFailed())

	var errorMsgs []string
	for _, err := range es {
		errorMsgs = append(errorMsgs, err.Error())
	}
	if len(errorMsgs) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMsgs, ", "))
	}
	return nil
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

func (x *XiaomiPushService) send(ctx context.Context, appID string, tokens []string, message *xp.Message) ([]string, error) {
	var newTokens []string
	var wg sync.WaitGroup
	client, ok := x.clients[appID]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentMeizuPushes <- struct{}{}
		wg.Add(1)
		x.status.AddXiaomiTotal(1)
		go func(notification *xp.Message, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentMeizuPushes
				wg.Done()
			}()

			fmt.Println("notification => ", notification)
			res, err := client.Send(ctx, notification, token)
			if err != nil || (res != nil && res.Code != 0) {
				if err == nil {
					err = errors.New(res.Reason)
				} else {
					es = append(es, err)
				}
				// 记录失败的 Token
				if res != nil && res.Code != 0 {
					newTokens = append(newTokens, token)
				}
				log.Printf("xiaomi send error: %s", err)
				x.status.AddXiaomiFailed(1)
			} else {
				log.Printf("xiaomi send success: %s", res.Reason)
				x.status.AddXiaomiSuccess(1)
			}
		}(message, token)
	}
	wg.Wait()
	if len(es) > 0 {
		var errorStrings []string
		for _, err := range es {
			errorStrings = append(errorStrings, err.Error())
		}
		allErrorsString := strings.Join(errorStrings, ", ")
		return nil, errors.New(allErrorsString)
	}
	return newTokens, nil
}

func (x *XiaomiPushService) SendMulticast(ctx context.Context, req interface{}, opt ...MulticastOption) error {
	//TODO implement me
	panic("implement me")
}

func (x *XiaomiPushService) Subscribe(ctx context.Context, req interface{}, opt ...SubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (x *XiaomiPushService) Unsubscribe(ctx context.Context, req interface{}, opt ...UnsubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (x *XiaomiPushService) SendToTopic(ctx context.Context, req interface{}, opt ...TopicOption) error {
	//TODO implement me
	panic("implement me")
}

func (x *XiaomiPushService) CheckDevice(ctx context.Context, req interface{}, opt ...CheckDeviceOption) bool {
	//TODO implement me
	panic("implement me")
}

func (x *XiaomiPushService) Name() string {
	//TODO implement me
	panic("implement me")
}
