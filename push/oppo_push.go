package push

import (
	"context"
	"errors"
	"fmt"
	op "github.com/316014408/oppo-push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/consts"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/status"
	"github.com/go-logr/logr"
	"log"
	"strings"
	"sync"
)

var (
	MaxConcurrentOppoPushes = make(chan struct{}, 100)
)

// OppoService 实现oppo推送，实现 PushService 接口
type OppoService struct {
	clients map[string]*op.OppoPush
	status  *status.StateStorage
	logger  logr.Logger
}

func NewOppoService(cfg *config.Config, logger logr.Logger) *OppoService {
	s := &OppoService{
		clients: map[string]*op.OppoPush{},
		status:  status.StatStorage,
		logger:  logger,
	}

	for _, v := range cfg.Oppo {
		if !v.Enabled || v.Enabled && (v.AppID == "" || v.AppKey == "" || v.AppSecret == "") {
			panic("push not enabled or misconfigured")
		}
		client := op.NewClient(v.AppKey, v.AppSecret)
		s.clients[v.AppID] = client
	}

	return s
}

func (o *OppoService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.OppoPushNotification)
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

	if retry > 0 && retry < maxRetry {
		maxRetry = retry
	}

	if err := o.checkNotification(req); err != nil {
		return err
	}

	notification, err := o.buildNotification(req)
	if err != nil {
		return err
	}

	if so.DryRun {
		return nil
	}

	for {
		newTokens, err := o.send(req.AppID, req.Tokens, notification)
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
		MaxConcurrentOppoPushes <- struct{}{}
		wg.Add(1)
		o.status.AddOppoTotal(1)
		go func(notification *op.Message, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentOppoPushes
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
				o.status.AddOppoFailed(1)
			} else {
				log.Printf("oppo send success: %s", res.Message)
				o.status.AddOppoSuccess(1)
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

func (o *OppoService) Name() string {
	return consts.PlatformOppo.String()
}
