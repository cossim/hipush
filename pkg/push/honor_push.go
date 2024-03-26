package push

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	hClient "github.com/cossim/hipush/pkg/client/push"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-logr/logr"
	"log"
	"net/http"
)

var (
	// MaxConcurrentHonorPushes pool to limit the number of concurrent iOS pushes
	MaxConcurrentHonorPushes = make(chan struct{}, 100)
)

// HonorService 荣耀推送，实现了 PushService 接口
type HonorService struct {
	clients        map[string]*hClient.HonorPushClient
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewHonorService(cfg *config.Config, logger logr.Logger) *HonorService {
	s := &HonorService{
		clients:        make(map[string]*hClient.HonorPushClient),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
	}

	for _, v := range cfg.Honor {
		if !v.Enabled {
			continue
		}
		if v.ClientID == "" || v.ClientSecret == "" {
			panic("push not enabled or misconfigured")
		}
		client := hClient.NewHonorPush(v.ClientID, v.ClientSecret)
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

func (h *HonorService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = h.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if err := h.checkNotification(req); err != nil {
		return nil, err
	}

	notification := h.buildAndroidNotification(req, so)

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return h.send(ctx, appid, token, notification)
	}

	resp, err := RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	taskid, err := h.getTaskIDFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: taskid}, nil
}

// getTaskIDFromResponse 从 Response 结构体中获取 task_id 字段
func (h *HonorService) getTaskIDFromResponse(response *Response) (string, error) {
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		return "", err
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(marshal, &dataMap); err != nil {
		return "", err
	}
	taskID, ok := dataMap["requestId"].(string)
	if !ok {
		return "", errors.New("task_id 字段不是 string 类型")
	}
	return taskID, nil
}

func (h *HonorService) GetTasksStatus(ctx context.Context, appid string, taskID []string, obj push.TaskObjectList) error {
	return nil
}

func (h *HonorService) send(ctx context.Context, appid string, token string, notification *hClient.SendMessageRequest) (*Response, error) {
	client, ok := h.clients[appid]
	if !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	h.status.AddHonorTotal(1)

	resp := &Response{Code: Fail}
	notification.Token = []string{token}
	res, err := client.SendMessage(ctx, appid, notification)
	if err != nil {
		log.Printf("honor send error: %s", err)
		h.status.AddHonorFailed(1)
		resp.Msg = err.Error()
	} else if res != nil && res.Code != http.StatusOK {
		if len(res.Data.ExpireTokens) > 0 {
			log.Printf("honor send expire tokens: %s", res.Data.ExpireTokens)
		}
		log.Printf("honor send error: %s", res.Message)
		h.status.AddHonorFailed(1)
		err = errors.New(res.Message)
		resp.Code = res.Code
		resp.Msg = res.Message
	} else {
		log.Printf("honor send success: %s", res.Message)
		h.status.AddHonorSuccess(1)
		resp.Code = Success
		resp.Msg = res.Message
		resp.Data = res
	}

	return resp, err

	//var es []error
	//
	//for _, token := range tokens {
	//	// occupy push slot
	//	MaxConcurrentHuaweiPushes <- struct{}{}
	//	wg.Add(1)
	//	h.status.AddHonorTotal(1)
	//	go func(notification *hClient.SendMessageRequest, token string) {
	//		defer func() {
	//			// free push slot
	//			<-MaxConcurrentHuaweiPushes
	//			wg.Done()
	//		}()
	//		res, err := client.SendMessage(ctx, appid, notification)
	//		if err != nil || (res != nil && res.Code != 200) {
	//			if err == nil {
	//				err = errors.New(res.Content)
	//			} else {
	//				es = append(es, err)
	//			}
	//			// 记录失败的 Token
	//			if res != nil && res.Code != 200 {
	//				newTokens = append(newTokens, res.Data.FailTokens...)
	//			}
	//			if len(res.Data.ExpireTokens) > 0 {
	//				log.Printf("honor send expire tokens: %s", res.Data.ExpireTokens)
	//			}
	//			log.Printf("honor send error: %s", err)
	//			h.status.AddHonorFailed(1)
	//
	//		} else {
	//			log.Printf("honor send success: %s", res.Content)
	//			h.status.AddHonorSuccess(1)
	//		}
	//	}(notification, token)
	//}
	//wg.Wait()
	//if len(es) > 0 {
	//	var errorStrings []string
	//	for _, err := range es {
	//		errorStrings = append(errorStrings, err.Error())
	//	}
	//	allErrorsString := strings.Join(errorStrings, ", ")
	//	return nil, errors.New(allErrorsString)
	//}
	//return newTokens, nil
}

func (h *HonorService) checkNotification(req push.SendRequest) error {
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

func (h *HonorService) buildAndroidNotification(req push.SendRequest, so *push.SendOptions) *hClient.SendMessageRequest {
	// 构建通知栏消息
	notification := &hClient.Notification{
		Title: req.GetTitle(),
		Body:  req.GetContent(),
		Image: req.GetIcon(),
	}

	// 构建 Android 平台的通知消息
	androidNotification := &hClient.AndroidNotification{
		Title: req.GetTitle(),
		Body:  req.GetContent(),
		Image: req.GetIcon(),
		//NotifyID: req.NotifyId,
		//Badge: &hClient.BadgeNotification{
		//	AddNum:     req.Badge.AddNum,
		//	SetNum:     req.Badge.SetNum,
		//	BadgeClass: req.Badge.BadgeClass,
		//},
		//ClickAction: &hClient.ClickAction{
		//	Type:   req.ClickAction.Action,
		//	Intent: req.ClickAction.Activity,
		//	URL:    req.ClickAction.Url,
		//	Action: req.ClickAction.Activity,
		//},
	}

	var targetUserType int

	if so.Development {
		targetUserType = 1
	}

	// 构建 Android 平台消息推送配置
	androidConfig := &hClient.AndroidConfig{
		//TTL:            req.GetTTL(),      // 设置消息缓存时间
		BiTag: "", // 设置批量任务消息标识
		//Data:           string(data), // 设置自定义消息负载
		Notification:   androidNotification,
		TargetUserType: targetUserType, // 设置目标用户类型
	}

	// 构建发送消息请求
	sendMessageReq := &hClient.SendMessageRequest{
		//Data:         string(data), // 设置自定义消息负载
		Notification: notification,
		Android:      androidConfig,
		Token:        req.GetToken(),
	}

	return sendMessageReq
}

func (h *HonorService) Name() string {
	return consts.PlatformHonor.String()
}
