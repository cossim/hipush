package push

import (
	"context"
	"errors"
	"fmt"
	op "github.com/316014408/oppo-push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/notify"
	"log"
	"strings"
	"sync"
)

// OppoService 实现vivo推送，必须实现PushService接口
type OppoService struct {
	clients map[string]*op.OppoPush
}

func NewOppoService(cfg *config.Config) (*OppoService, error) {
	s := &OppoService{
		clients: map[string]*op.OppoPush{},
	}

	for _, v := range cfg.Oppo {
		if !v.Enabled || v.Enabled && (v.AppID == "" || v.AppKey == "" || v.AppSecret == "") {
			return nil, errors.New("push not enabled or misconfigured")
		}
		client := op.NewClient(v.AppKey, v.AppSecret)
		s.clients[v.AppID] = client
	}

	return s, nil
}

func (o *OppoService) Send(ctx context.Context, request interface{}) error {
	req, ok := request.(*notify.OppoPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	var (
		retry      = req.Option.Retry
		maxRetry   = retry
		retryCount = 0
	)

	// 重试计数
	if maxRetry <= 0 {
		maxRetry = DefaultMaxRetry // 设置一个默认的最大重试次数
	}
	if retry > 0 && retry < maxRetry {
		maxRetry = retry
	}

	var es []error

	if err := o.checkNotification(req); err != nil {
		return err
	}

	notification, err := o.buildNotification(req)
	if err != nil {
		return err
	}

	if req.Option.DryRun {
		return nil
	}

	for {
		newTokens, err := o.send(req.AppID, req.Tokens, notification)
		if err != nil {
			log.Printf("sendNotifications error => %v", err)
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

	var errorMsgs []string
	for _, err := range es {
		errorMsgs = append(errorMsgs, err.Error())
	}
	if len(errorMsgs) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMsgs, ", "))
	}
	return nil
}

func (o *OppoService) send(appID string, tokens []string, message *op.Message) ([]string, error) {
	var newTokens []string
	var wg sync.WaitGroup
	fmt.Println("appID:", appID)
	client, ok := o.clients[appID]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentIOSPushes <- struct{}{}
		wg.Add(1)
		go func(notification *op.Message, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentIOSPushes
				wg.Done()
			}()

			notification.SetTargetValue(token)
			fmt.Println("notification => ", notification.String())
			res, err := client.Unicast(notification)
			if err != nil || (res != nil && res.Code != 0) {
				if err == nil {
					err = errors.New(res.Message)
				} else {
					es = append(es, err)
				}
				// 记录失败的 Token
				if res != nil && res.Code != 0 {
					newTokens = append(newTokens, token)
				}
				log.Printf("oppo send error: %s", err)
			} else {
				log.Printf("oppo send success: %s", res.Message)
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

func (o *OppoService) checkNotification(req *notify.OppoPushNotification) error {
	if len(req.Tokens) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.Title == "" {
		return errors.New("title cannot be empty")
	}

	if req.Message == "" {
		return errors.New("message cannot be empty")
	}

	return nil
}

func (o *OppoService) buildNotification(req *notify.OppoPushNotification) (*op.Message, error) {
	m := op.NewMessage(req.Title, req.Message).
		SetSubTitle(req.Subtitle).
		SetTargetType(2)
	if req.ClickAction != nil {
		if req.ClickAction.Action == 1 || req.ClickAction.Action == 4 {
			m.SetClickActionActivity(req.ClickAction.Activity)
		}
		if req.ClickAction.Action == 2 {
			m.SetClickActionUrl(req.ClickAction.Url)
		}
		m.SetClickActionType(req.ClickAction.Action)
		m.SetActionParameters(req.ClickAction.Parameters)
	}

	return m, nil
}

func (o *OppoService) MulticastSend(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (o *OppoService) Subscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (o *OppoService) Unsubscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (o *OppoService) SendToTopic(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (o *OppoService) SendToCondition(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (o *OppoService) CheckDevice(ctx context.Context, req interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (o *OppoService) Name() string {
	//TODO implement me
	panic("implement me")
}
