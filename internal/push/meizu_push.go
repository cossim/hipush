package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mzp "github.com/cossim/go-meizu-push-sdk"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/notify"
	"log"
	"strings"
	"sync"
)

var (
	MaxConcurrentMeizuPushes = make(chan struct{}, 100)
)

// MeizuService 实现vivo推送，必须实现PushService接口
type MeizuService struct {
	clients map[string]func(token, message string) mzp.PushResponse
}

func NewMeizuService(cfg *config.Config) (*MeizuService, error) {
	s := &MeizuService{
		clients: make(map[string]func(token, message string) mzp.PushResponse),
	}

	for _, v := range cfg.Meizu {
		if !v.Enabled || v.Enabled && (v.AppID == "" || v.AppKey == "") {
			return nil, errors.New("push not enabled or misconfigured")
		}
		s.clients[v.AppID] = func(token, message string) mzp.PushResponse {
			appid := v.AppID
			appkey := v.AppKey
			return mzp.PushNotificationMessageByPushId(appid, token, message, appkey)
		}
	}

	return s, nil
}

func (m *MeizuService) Send(ctx context.Context, request interface{}, opt SendOption) error {
	req, ok := request.(*notify.MeizuPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	var (
		retry      = opt.Retry
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

	if err := m.checkNotification(req); err != nil {
		return err
	}

	notification, err := m.buildNotification(req)
	if err != nil {
		return err
	}

	if opt.DryRun {
		return nil
	}

	for {
		newTokens, err := m.send(req.AppID, req.Tokens, notification)
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

func (m *MeizuService) send(appid string, tokens []string, message string) ([]string, error) {
	var newTokens []string
	var wg sync.WaitGroup
	pushFunc, ok := m.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentXiaomiPushes <- struct{}{}
		wg.Add(1)
		go func(notification string, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentXiaomiPushes
				wg.Done()
			}()

			fmt.Println("notification => ", notification)
			res := pushFunc(token, notification)
			if res.GetCode() != 200 {
				es = append(es, errors.New(res.GetMessage()))
				log.Printf("oppo send error: %s", res.GetMessage())
				// 记录失败的 Token
				newTokens = append(newTokens, token)
			} else {
				log.Printf("oppo send success: %s", res.GetMessage())
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

func (m *MeizuService) checkNotification(req *notify.MeizuPushNotification) error {
	if len(req.Tokens) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.Title == "" {
		return errors.New("title cannot be empty")
	}

	if req.Content == "" {
		return errors.New("message cannot be empty")
	}

	if req.IsScheduled && (req.ScheduledStartTime == "" || req.ScheduledEndTime == "") {
		return errors.New("scheduled time cannot be empty")
	}

	return nil
}

func (m *MeizuService) buildNotification(req *notify.MeizuPushNotification) (string, error) {
	msg := mzp.BuildNotificationMessage()
	msg.NoticeBarInfo.Title = req.Title
	msg.NoticeBarInfo.Content = req.Content
	msg.ClickTypeInfo = mzp.ClickTypeInfo{
		ClickType:  req.ClickAction.Action,
		Url:        req.ClickAction.Url,
		Parameters: req.ClickAction.Parameters,
		Activity:   req.ClickAction.Activity,
	}

	offLine := 0
	if req.OffLine {
		offLine = 1
	}
	msg.PushTimeInfo = mzp.PushTimeInfo{
		OffLine:   offLine,
		ValidTime: req.TTL,
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(message), nil
}

func (m *MeizuService) MulticastSend(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MeizuService) Subscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MeizuService) Unsubscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MeizuService) SendToTopic(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MeizuService) SendToCondition(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MeizuService) CheckDevice(ctx context.Context, req interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MeizuService) Name() string {
	//TODO implement me
	panic("implement me")
}
