package push

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-logr/logr"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"log"
	"net"
	"path/filepath"
	"time"
)

var (
	_ push.PushService = &APNsService{}

	idleConnTimeout = 90 * time.Second
	tlsDialTimeout  = 20 * time.Second
	tcpKeepAlive    = 60 * time.Second
)

const (
	dotP8  = ".p8"
	dotPEM = ".pem"
	dotP12 = ".p12"
)

// APNsService 实现APNs推送，实现 PushService 接口
type APNsService struct {
	clients        map[string]*apns2.Client
	appNameToIDMap map[string]string
	status         *status.StateStorage
	logger         logr.Logger
}

func NewAPNsService(cfg *config.Config, logger logr.Logger) *APNsService {
	var ext string
	var err error
	var authKey *ecdsa.PrivateKey
	//var certificateKey tls.Certificate

	s := &APNsService{
		clients:        make(map[string]*apns2.Client),
		appNameToIDMap: make(map[string]string),
		status:         status.StatStorage,
		logger:         logger,
	}

	for _, v := range cfg.IOS {
		if !v.Enabled {
			continue
		}

		if v.KeyPath == "" || v.AppID == "" {
			msg := "you should provide ios.KeyPath and ios.AppID"
			panic(msg)
		}

		ext = filepath.Ext(v.KeyPath)
		switch ext {
		case dotP12:
			//certificateKey, err = certificate.FromP12File(v.KeyPath, v.Password)
		case dotPEM:
			//certificateKey, err = certificate.FromPemFile(v.KeyPath, v.Password)
		case dotP8:
			authKey, err = token.AuthKeyFromFile(v.KeyPath)
			if v.KeyID == "" || v.TeamID == "" || v.AppID == "" {
				msg := "you should provide ios.KeyID and ios.TeamID for p8 token"
				panic(msg)
			}
			token := &token.Token{
				AuthKey: authKey,
				// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
				KeyID: v.KeyID,
				// TeamID from developer account (View Account -> Membership)
				TeamID: v.TeamID,
			}
			client, err := s.newApnsTokenClient(v.Production, token)
			if err != nil {
				panic(err)
			}
			s.clients[v.AppID] = client
			if v.AppName != "" {
				s.appNameToIDMap[v.AppName] = v.AppID
			}
		default:
			err = errors.New("wrong certificate key extension")
		}
		if err != nil {
			panic(err)
		}
	}

	return s
}

func (a *APNsService) newApnsTokenClient(production bool, token *token.Token) (*apns2.Client, error) {
	var client *apns2.Client

	if production {
		client = apns2.NewTokenClient(token).Production()
	} else {
		client = apns2.NewTokenClient(token).Development()
	}

	//transport := &http.Transport{
	//	DialTLS:         DialTLS(nil),
	//	Proxy:           http.DefaultTransport.(*http.Transport).Proxy,
	//	IdleConnTimeout: idleConnTimeout,
	//}
	//
	//h2Transport, err := http2.ConfigureTransports(transport)
	//if err != nil {
	//	return nil, err
	//}
	//
	//h2Transport.ReadIdleTimeout = 1 * time.Second
	//h2Transport.PingTimeout = 1 * time.Second
	//
	//client.HTTPClient.Transport = transport
	return client, nil
}

// DialTLS is the default dial function for creating TLS connections for
// non-proxied HTTPS requests.
var DialTLS = func(cfg *tls.Config) func(network, addr string) (net.Conn, error) {
	return func(network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   tlsDialTimeout,
			KeepAlive: tcpKeepAlive,
		}
		return tls.DialWithDialer(dialer, network, addr, cfg)
	}
}

func (a *APNsService) Send(ctx context.Context, req push.SendRequest, opt ...push.SendOption) (*push.SendResponse, error) {
	so := &push.SendOptions{}
	so.ApplyOptions(opt)

	var appid string
	var ok bool
	if req.GetAppID() != "" {
		appid = req.GetAppID()
	} else if req.GetAppName() != "" {
		appid, ok = a.appNameToIDMap[req.GetAppName()]
		if !ok {
			return nil, ErrInvalidAppID
		}
	} else {
		return nil, ErrInvalidAppID
	}

	if err := a.checkNotification(req); err != nil {
		return nil, err
	}

	notification, err := a.buildNotification(req)
	if err != nil {
		return nil, err
	}

	if so.DryRun {
		return nil, nil
	}

	send := func(ctx context.Context, token string) (*Response, error) {
		return a.send(appid, token, notification)
	}

	resp, err := RetrySend(ctx, send, req.GetToken(), so.Retry, so.RetryInterval, 100)
	if err != nil {
		return nil, err
	}

	taskid, err := a.getTaskIDFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return &push.SendResponse{TaskId: taskid}, nil
}

// getTaskIDFromResponse 从 Response 结构体中获取 RequestId
func (a *APNsService) getTaskIDFromResponse(response *Response) (string, error) {
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		return "", err
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(marshal, &dataMap); err != nil {
		return "", err
	}
	requestID, ok := dataMap["ApnsID"].(string)
	if !ok {
		return "", errors.New("RequestId 字段不是 string 类型")
	}
	return requestID, nil
}

func (a *APNsService) GetTasksStatus(ctx context.Context, appid string, taskID []string, obj push.TaskObjectList) error {
	return nil
}

// Sound sets the aps sound on the payload.
// https://developer.apple.com/documentation/usernotifications/generating-a-remote-notification#:~:text=%E8%AD%A6%E6%8A%A5%E7%9A%84%E5%A3%B0%E9%9F%B3%E3%80%82-,sound%E8%A1%A8%203.%E5%AD%97%E5%85%B8%E4%B8%AD%E5%8C%85%E5%90%AB%E7%9A%84%E9%94%AE,-%E9%92%A5%E5%8C%99
type Sound struct {
	// Critical 指示声音是否被标记为关键声音。关键声音通常用于需要立即引起用户注意的通知。
	// 值为 1 表示是关键声音，值为 0 表示不是关键声音。默认为 0。
	Critical int `json:"critical,omitempty"`
	// Name 声音的名称或标识符。
	// 通常是声音文件的名称，表示要播放的声音文件。默认为空字符串。
	Name string `json:"name,omitempty"`
	// Volume 声音的音量级别。
	// 值范围为 0.0 到 1.0，表示音量的相对级别。默认为 1.0。
	Volume float32 `json:"volume,omitempty"`
}

func (a *APNsService) buildNotification(req push.SendRequest) (*apns2.Notification, error) {
	topic := req.GetTopic()
	if topic == "" {
		topic = req.GetAppID()
	}
	notification := &apns2.Notification{
		ApnsID:     req.GetMessageID(),
		Topic:      topic,
		CollapseID: req.GetCollapseID(),
	}

	if req.GetTTL() != 0 {
		notification.Expiration = time.Unix(req.GetTTL(), 0)
	}

	if req.GetPriority() == "normal" {
		notification.Priority = apns2.PriorityLow
	} else if req.GetPriority() == "high" {
		notification.Priority = apns2.PriorityHigh
	}

	//if len(req.PushType) > 0 {
	//	notification.PushType = apns2.EPushType(req.PushType)
	//}

	payload := payload.NewPayload()

	// add alert object if message length > 0 and title is empty
	if len(req.GetContent()) > 0 && req.GetTitle() == "" {
		payload.Alert(req.GetContent())
	}

	// zero value for clear the badge on the app icon.
	//if req.Badge != nil && *req.Badge >= 0 {
	//	payload.Badge(*req.Badge)
	//}

	if req.GetMutableContent() {
		payload.MutableContent()
	}

	//switch req.GetSound().(type) {
	//// from http request binding
	//case map[string]interface{}:
	//	result := &Sound{}
	//	_ = mapstructure.Decode(req.GetSound(), &result)
	//	payload.Sound(result)
	//// from http request binding for non critical alerts
	//case string:
	//	payload.Sound(req.GetSound())
	//case Sound:
	//	payload.Sound(req.GetSound())
	//}

	//if len(req.SoundName) > 0 {
	//	payload.SoundName(req.SoundName)
	//}
	//
	//if req.SoundVolume > 0 {
	//	payload.SoundVolume(req.SoundVolume)
	//}

	if req.GetContentAvailable() {
		payload.ContentAvailable()
	}

	//if len(req.URLArgs) > 0 {
	//	payload.URLArgs(req.URLArgs)
	//}

	//if len(req.ThreadID) > 0 {
	//	payload.ThreadID(req.ThreadID)
	//}

	//for k, v := range req.GetData() {
	//	payload.Custom(k, v)
	//}

	payload.AlertTitle(req.GetTitle())
	payload.AlertBody(req.GetContent())
	payload.Category(req.GetCategory())

	//payload = iosAlertDictionary(payload, req)

	notification.Payload = payload

	return notification, nil
	//return notify.GetIOSNotification(req), nil
}

//func iosAlertDictionary(notificationPayload *payload.Payload, req push.SendRequest) *payload.Payload {
//	// Alert dictionary
//
//	if len(req.GetTitle()) > 0 {
//		notificationPayload.AlertTitle(req.GetTitle())
//	}
//
//	if len(req.InterruptionLevel) > 0 {
//		notificationPayload.InterruptionLevel(payload.EInterruptionLevel(req.InterruptionLevel))
//	}
//
//	if len(req.GetContent()) > 0 && len(req.GetTitle()) > 0 {
//		notificationPayload.AlertBody(req.GetContent())
//	}
//
//	if len(req.Alert.Title) > 0 {
//		notificationPayload.AlertTitle(req.Alert.Title)
//	}
//
//	// Apple Watch & Safari display this string as part of the notification interface.
//	if len(req.Alert.Subtitle) > 0 {
//		notificationPayload.AlertSubtitle(req.Alert.Subtitle)
//	}
//
//	if len(req.Alert.TitleLocKey) > 0 {
//		notificationPayload.AlertTitleLocKey(req.Alert.TitleLocKey)
//	}
//
//	if len(req.Alert.LocArgs) > 0 {
//		notificationPayload.AlertLocArgs(req.Alert.LocArgs)
//	}
//
//	if len(req.Alert.TitleLocArgs) > 0 {
//		notificationPayload.AlertTitleLocArgs(req.Alert.TitleLocArgs)
//	}
//
//	if len(req.Alert.Body) > 0 {
//		notificationPayload.AlertBody(req.Alert.Body)
//	}
//
//	if len(req.Alert.LaunchImage) > 0 {
//		notificationPayload.AlertLaunchImage(req.Alert.LaunchImage)
//	}
//
//	if len(req.Alert.LocKey) > 0 {
//		notificationPayload.AlertLocKey(req.Alert.LocKey)
//	}
//
//	if len(req.Alert.Action) > 0 {
//		notificationPayload.AlertAction(req.Alert.Action)
//	}
//
//	if len(req.Alert.ActionLocKey) > 0 {
//		notificationPayload.AlertActionLocKey(req.Alert.ActionLocKey)
//	}
//
//	// General
//	if len(req.Category) > 0 {
//		notificationPayload.Category(req.Category)
//	}
//
//	if len(req.Alert.SummaryArg) > 0 {
//		notificationPayload.AlertSummaryArg(req.Alert.SummaryArg)
//	}
//	if req.Alert.SummaryArgCount > 0 {
//		notificationPayload.AlertSummaryArgCount(req.Alert.SummaryArgCount)
//	}
//
//	return notificationPayload
//}

func (a *APNsService) send(appid string, token string, notification *apns2.Notification) (*Response, error) {
	if _, ok := a.clients[appid]; !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}
	resp := &Response{Code: Fail}
	a.status.AddIosTotal(1)
	notification.DeviceToken = token
	res, err := a.clients[appid].Push(notification)
	if err != nil {
		log.Printf("apns send error: %s", err)
		a.status.AddIosFailed(1)
		resp.Msg = err.Error()
	} else if res != nil && res.StatusCode != Success {
		log.Printf("apns send error: %s", res.Reason)
		a.status.AddIosFailed(1)
		err = errors.New(res.Reason)
		resp.Msg = res.Reason
	} else {
		log.Printf("apns send success: %v", res)
		a.status.AddIosSuccess(1)
		resp.Code = Success
		resp.Msg = res.Reason
		resp.Data = res
	}

	return resp, err
}

func (a *APNsService) checkNotification(req push.SendRequest) error {
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

func (a *APNsService) Name() string {
	return consts.PlatformIOS.String()
}
