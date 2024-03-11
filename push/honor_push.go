package push

import (
	"context"
	"errors"
	"fmt"
	hClient "github.com/cossim/hipush/client/push"
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
	// MaxConcurrentHonorPushes pool to limit the number of concurrent iOS pushes
	MaxConcurrentHonorPushes = make(chan struct{}, 100)
)

// HonorService 荣耀推送，实现了 PushService 接口
type HonorService struct {
	clients map[string]*hClient.HonorPushClient
	status  *status.StateStorage
	logger  logr.Logger
}

func NewHonorService(cfg *config.Config, logger logr.Logger) (*HonorService, error) {
	s := &HonorService{
		clients: make(map[string]*hClient.HonorPushClient),
		status:  status.StatStorage,
		logger:  logger,
	}

	for _, v := range cfg.Honor {
		if !v.Enabled {
			continue
		}
		if v.Enabled && (v.AppID == "" || v.ClientID == "" || v.ClientSecret == "") {
			return nil, errors.New("push not enabled or misconfigured")
		}
		s.clients[v.AppID] = hClient.NewHonorPush(v.ClientID, v.ClientSecret)
	}

	return s, nil
}

func (h *HonorService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.HonorPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	so := &SendOptions{}
	so.ApplyOptions(opt)

	var maxRetry = so.Retry
	var retryCount int
	var es []error

	if err := h.checkNotification(req); err != nil {
		return err
	}

	notification := h.buildAndroidNotification(req)

	if so.DryRun {
		return nil
	}

	for {
		newTokens, err := h.send(ctx, req.AppID, req.Tokens, notification)
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

func (h *HonorService) send(ctx context.Context, appid string, tokens []string, notification *hClient.SendMessageRequest) ([]string, error) {
	var newTokens []string
	var wg sync.WaitGroup

	client, ok := h.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentHuaweiPushes <- struct{}{}
		wg.Add(1)
		h.status.AddHonorTotal(1)
		go func(notification *hClient.SendMessageRequest, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentHuaweiPushes
				wg.Done()
			}()
			res, err := client.SendMessage(ctx, appid, notification)
			if err != nil || (res != nil && res.Code != 200) {
				if err == nil {
					err = errors.New(res.Message)
				} else {
					es = append(es, err)
				}
				// 记录失败的 Token
				if res != nil && res.Code != 200 {
					newTokens = append(newTokens, res.Data.FailTokens...)
				}
				if len(res.Data.ExpireTokens) > 0 {
					log.Printf("honor send expire tokens: %s", res.Data.ExpireTokens)
				}
				log.Printf("honor send error: %s", err)
				h.status.AddHonorFailed(1)

			} else {
				log.Printf("honor send success: %s", res.Message)
				h.status.AddHonorSuccess(1)
			}
		}(notification, token)
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

func (h *HonorService) checkNotification(req *notify.HonorPushNotification) error {
	if len(req.Tokens) == 0 {
		return errors.New("tokens cannot be empty")
	}

	if req.Title == "" {
		return errors.New("title cannot be empty")
	}

	if req.Content == "" {
		return errors.New("content cannot be empty")
	}

	return nil
}

func (h *HonorService) buildAndroidNotification(req *notify.HonorPushNotification) *hClient.SendMessageRequest {
	// 构建通知栏消息
	notification := &hClient.Notification{
		Title: req.Title,
		Body:  req.Content,
	}

	// 构建 Android 平台的通知消息
	androidNotification := &hClient.AndroidNotification{
		Title: req.Title,
		Body:  req.Content,
		Image: req.Image,
		ClickAction: &hClient.ClickAction{
			Type:   req.ClickAction.Action,
			Intent: req.ClickAction.Activity,
			URL:    req.ClickAction.Url,
			Action: req.ClickAction.Activity,
		},
	}

	var targetUserType int

	if req.Development {
		targetUserType = 1
	}

	// 构建 Android 平台消息推送配置
	androidConfig := &hClient.AndroidConfig{
		TTL:            req.TTL,  // 设置消息缓存时间
		BiTag:          "",       // 设置批量任务消息标识
		Data:           req.Data, // 设置自定义消息负载
		Notification:   androidNotification,
		TargetUserType: targetUserType, // 设置目标用户类型
	}

	// 构建发送消息请求
	sendMessageReq := &hClient.SendMessageRequest{
		Data:         req.Data, // 设置自定义消息负载
		Notification: notification,
		Android:      androidConfig,
		Token:        req.Tokens,
	}

	return sendMessageReq
}

func (h *HonorService) Name() string {
	return consts.PlatformHonor.String()
}
