package push

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/status"
	"github.com/go-logr/logr"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	// MaxConcurrentIOSPushes pool to limit the number of concurrent iOS pushes
	MaxConcurrentIOSPushes chan struct{}

	idleConnTimeout = 90 * time.Second
	tlsDialTimeout  = 20 * time.Second
	tcpKeepAlive    = 60 * time.Second

	doOnce sync.Once
)

const (
	dotP8  = ".p8"
	dotPEM = ".pem"
	dotP12 = ".p12"
)

// APNsService 实现APNs推送，实现 PushService 接口
type APNsService struct {
	clients map[string]*apns2.Client
	status  *status.StateStorage
	logger  logr.Logger
}

func NewAPNsService(cfg *config.Config, logger logr.Logger) (*APNsService, error) {
	var ext string
	var err error
	var authKey *ecdsa.PrivateKey
	//var certificateKey tls.Certificate

	s := &APNsService{
		clients: make(map[string]*apns2.Client),
		status:  status.StatStorage,
		logger:  logger,
	}

	for _, v := range cfg.IOS {
		if v.KeyPath != "" && v.Enabled {
			ext = filepath.Ext(v.KeyPath)
			switch ext {
			case dotP12:
				//certificateKey, err = certificate.FromP12File(v.KeyPath, v.Password)
			case dotPEM:
				//certificateKey, err = certificate.FromPemFile(v.KeyPath, v.Password)
			case dotP8:
				authKey, err = token.AuthKeyFromFile(v.KeyPath)
				if v.KeyID == "" || v.TeamID == "" {
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
					return nil, err
				}
				s.clients[v.AppID] = client
			default:
				err = errors.New("wrong certificate key extension")
			}
			if err != nil {
				return nil, err
			}
		}
	}

	doOnce.Do(func() {
		MaxConcurrentIOSPushes = make(chan struct{}, 100)
	})

	return s, nil
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

func (a *APNsService) Send(ctx context.Context, request interface{}, opt ...SendOption) error {
	req, ok := request.(*notify.ApnsPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	so := &SendOptions{}
	so.ApplyOptions(opt)

	var (
		retry      = so.Retry
		maxRetry   = so.Retry
		retryCount = 0
		es         []error
	)

	if retry > 0 && retry < maxRetry {
		maxRetry = retry
	}

	if err := a.checkNotification(req); err != nil {
		return err
	}

	notification, err := a.buildNotification(req)
	if err != nil {
		return err
	}

	if so.DryRun {
		return nil
	}

	for {
		newTokens, err := a.send(req.AppID, req.Tokens, notification)
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

func (a *APNsService) MapPushRequestToApnsPushNotification(req PushRequest) (*notify.ApnsPushNotification, error) {
	badge := req.GetBadge()

	alert, err := a.MapToAlert(req.GetAlert())
	if err != nil {
		log.Println("MapToAlert error => ", err)
		return nil, err
	}

	apnsNotification := &notify.ApnsPushNotification{
		Tokens:            req.GetTokens(),
		Priority:          req.GetPriority(),
		Title:             req.GetTitle(),
		Content:           req.GetMessage(),
		Expiration:        req.GetExpiration(),
		ApnsID:            req.GetApnsID(),
		CollapseID:        req.GetCollapseID(),
		Topic:             req.GetTopic(),
		PushType:          req.GetPushType(),
		Badge:             &badge,
		Category:          req.GetCategory(),
		ThreadID:          req.GetThreadID(),
		URLArgs:           req.GetURLArgs(),
		Alert:             alert,
		ContentAvailable:  req.IsContentAvailable(),
		MutableContent:    req.IsMutableContent(),
		Production:        !req.IsDevelopment(),
		Development:       req.IsDevelopment(),
		SoundName:         req.GetSoundName(),
		SoundVolume:       req.GetSoundVolume(),
		Apns:              req.GetApns(),
		InterruptionLevel: req.GetInterruptionLevel(),
		Sound:             req.GetSound(),
		Data:              req.GetData(),
	}
	return apnsNotification, nil
}

// MapToAlert 将接口类型的数据转换为 notify.Alert 结构体
func (a *APNsService) MapToAlert(data interface{}) (notify.Alert, error) {
	alertMap, ok := data.(map[string]interface{})
	if !ok {
		return notify.Alert{}, errors.New("data is not in the expected format")
	}

	var toString = func(value interface{}) string {
		if str, ok := value.(string); ok {
			return str
		}
		return ""
	}

	var toStringSlice = func(value interface{}) []string {
		if slice, ok := value.([]interface{}); ok {
			strSlice := make([]string, len(slice))
			for i, v := range slice {
				if str, ok := v.(string); ok {
					strSlice[i] = str
				}
			}
			return strSlice
		}
		return nil
	}

	var toInt = func(value interface{}) int {
		if intValue, ok := value.(int); ok {
			return intValue
		}
		return 0
	}

	alert := notify.Alert{
		Action:          toString(alertMap["action"]),
		ActionLocKey:    toString(alertMap["action-loc-key"]),
		Body:            toString(alertMap["body"]),
		LaunchImage:     toString(alertMap["launch-image"]),
		LocArgs:         toStringSlice(alertMap["loc-args"]),
		LocKey:          toString(alertMap["loc-key"]),
		Title:           toString(alertMap["title"]),
		Subtitle:        toString(alertMap["subtitle"]),
		TitleLocArgs:    toStringSlice(alertMap["title-loc-args"]),
		TitleLocKey:     toString(alertMap["title-loc-key"]),
		SummaryArg:      toString(alertMap["summary-arg"]),
		SummaryArgCount: toInt(alertMap["summary-arg-count"]),
	}
	return alert, nil
}

func (a *APNsService) buildNotification(req *notify.ApnsPushNotification) (*apns2.Notification, error) {
	return notify.GetIOSNotification(req), nil
}

func (a *APNsService) send(appid string, tokens []string, notification *apns2.Notification) ([]string, error) {
	var newTokens []string
	var wg sync.WaitGroup

	if _, ok := a.clients[appid]; !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentIOSPushes <- struct{}{}
		wg.Add(1)
		a.status.AddIosTotal(1)
		go func(notification *apns2.Notification, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentIOSPushes
				wg.Done()
			}()

			notification.DeviceToken = token
			res, err := a.clients[appid].Push(notification)
			if err != nil || (res != nil && res.StatusCode != http.StatusOK) {
				if err == nil {
					err = errors.New(res.Reason)
				} else {
					es = append(es, err)
				}
				// 记录失败的 Token
				if res != nil && res.StatusCode >= http.StatusInternalServerError {
					newTokens = append(newTokens, token)
				}
				log.Printf("apns send error: %s", err)
				a.status.AddIosFailed(1)
			} else {
				log.Printf("apns send success: %s", res.Reason)
				a.status.AddIosSuccess(1)
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

func (a *APNsService) checkNotification(req *notify.ApnsPushNotification) error {
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

func (a *APNsService) SendMulticast(ctx context.Context, req interface{}, opt ...MulticastOption) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Subscribe(ctx context.Context, req interface{}, opt ...SubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Unsubscribe(ctx context.Context, req interface{}, opt ...UnsubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) SendToTopic(ctx context.Context, req interface{}, opt ...TopicOption) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) CheckDevice(ctx context.Context, req interface{}, opt ...CheckDeviceOption) bool {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Name() string {
	//TODO implement me
	panic("implement me")
}
