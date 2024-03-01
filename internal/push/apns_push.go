package push

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/notify"
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

// APNsService 实现APNs推送
type APNsService struct {
	clients map[string]*apns2.Client
}

func NewAPNsService(cfg *config.Config) (*APNsService, error) {
	var ext string
	var err error
	var authKey *ecdsa.PrivateKey
	//var certificateKey tls.Certificate

	s := &APNsService{
		clients: make(map[string]*apns2.Client),
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

func (a *APNsService) Send(ctx context.Context, request interface{}) error {
	req, ok := request.(*notify.ApnsPushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	marshal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("req => ", string(marshal))

	retry := req.Retry

	var maxRetry = retry
	if maxRetry <= 0 {
		maxRetry = DefaultMaxRetry // 设置一个默认的最大重试次数
	}
	if retry > 0 && retry < maxRetry {
		maxRetry = retry
	}

	// 重试计数
	retryCount := 0

	//notification, err := a.MapPushRequestToApnsPushNotification(req)
	//if err != nil {
	//	return err
	//}

	var es error

	for {
		newTokens, err := a.sendNotifications(req)
		if err != nil {
			log.Printf("sendNotifications error => %v", err)
			es = err
		}
		// 如果有重试的 Token，并且未达到最大重试次数，则进行重试
		if len(newTokens) > 0 && retryCount < maxRetry {
			retryCount++
			req.Tokens = newTokens
		} else {
			break
		}
	}

	return es
}

func (a *APNsService) MapPushRequestToApnsPushNotification(req PushRequest) (*notify.ApnsPushNotification, error) {
	badge := req.GetBadge()

	alert, err := a.MapToAlert(req.GetAlert())
	if err != nil {
		log.Println("MapToAlert error => ", err)
		return nil, err
	}

	apnsNotification := &notify.ApnsPushNotification{
		Retry:             req.GetRetry(),
		Tokens:            req.GetTokens(),
		Priority:          req.GetPriority(),
		Title:             req.GetTitle(),
		Message:           req.GetMessage(),
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

func (a *APNsService) sendNotifications(req *notify.ApnsPushNotification) ([]string, error) {
	var newTokens []string
	notification := notify.GetIOSNotification(req)
	var wg sync.WaitGroup

	if _, ok := a.clients[req.ApnsID]; !ok {
		return nil, errors.New("invalid appid or appid push is not enabled")
	}

	var es []error

	for _, token := range req.Tokens {
		// occupy push slot
		MaxConcurrentIOSPushes <- struct{}{}
		wg.Add(1)
		go func(notification apns2.Notification, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentIOSPushes
				wg.Done()
			}()

			notification.DeviceToken = token
			res, err := a.clients[req.ApnsID].Push(&notification)
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
			} else {
				log.Printf("apns send success: %s", res.Reason)
			}
		}(*notification, token)
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

func (a *APNsService) MulticastSend(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Subscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Unsubscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) SendToTopic(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) SendToCondition(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) CheckDevice(ctx context.Context, req interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Name() string {
	//TODO implement me
	panic("implement me")
}
