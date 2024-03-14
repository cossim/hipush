package push

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidAppID = errors.New("invalid appid or appid push is not enabled")
)

// PushRequest 获取不同厂商的推送请求参数
type PushRequest interface {
	Get() interface{}

	// GetRetry 获取推送消息的重试次数
	GetRetry() int

	// GetTokens 获取要发送消息的设备 token 列表
	GetTokens() []string

	// GetTitle 获取推送消息的标题
	GetTitle() string

	// GetMessage 获取推送消息的内容
	GetMessage() string

	// GetTopic 获取推送消息的主题（仅用于特定平台，如 APNs）
	GetTopic() string

	// GetKey 获取推送消息的密钥（仅用于特定平台，如 APNs）
	GetKey() string

	// GetCategory 获取推送消息的类别（仅用于特定平台，如 APNs）
	GetCategory() string

	// GetSound 获取推送消息的声音
	GetSound() interface{}

	//GetAlert 获取推送消息的警报信息（仅用于特定平台，如 APNs）
	GetAlert() interface{}

	// GetBadge 获取推送消息的应用图标角标数
	GetBadge() int

	// GetThreadID 获取推送消息的线程 ID（仅用于特定平台，如 APNs）
	GetThreadID() string

	// GetData 获取推送消息的自定义数据
	GetData() map[string]interface{}

	// GetImage 获取推送消息的图片（仅用于特定平台，如 APNs）
	GetImage() string

	// GetID 获取推送消息的唯一标识符（仅用于特定平台，如 APNs）
	GetID() string

	// GetPushType 获取推送消息的类型（仅用于特定平台，如 APNs）
	GetPushType() string

	// GetPriority 获取推送消息的优先级（仅用于特定平台，如 APNs）
	GetPriority() string

	// IsContentAvailable 检查推送消息是否可用于内容可用性（仅用于特定平台，如 APNs）
	IsContentAvailable() bool

	// IsMutableContent 检查推送消息是否可变（仅用于特定平台，如 APNs）
	IsMutableContent() bool

	// IsDevelopment 检查推送消息是否用于开发环境（仅用于特定平台，如 APNs）
	IsDevelopment() bool

	// GetExpiration 获取推送消息的过期时间
	GetExpiration() *int64

	// GetApnsID 获取推送消息的唯一标识符（仅用于特定平台，如 APNs）
	GetApnsID() string

	// GetCollapseID 获取推送消息的折叠 ID（仅用于特定平台，如 APNs）
	GetCollapseID() string

	// GetURLArgs 获取推送消息的 URL 参数（仅用于特定平台，如 APNs）
	GetURLArgs() []string

	// GetSoundName 获取推送消息的声音名称（仅用于特定平台，如 APNs）
	GetSoundName() string

	// GetSoundVolume 获取推送消息的声音音量（仅用于特定平台，如 APNs）
	GetSoundVolume() float32

	// GetApns 获取推送消息的 APNs 通知（仅用于特定平台，如 APNs）
	GetApns() map[string]interface{}

	// GetInterruptionLevel 获取推送消息的中断等级（仅用于特定平台，如 APNs）
	GetInterruptionLevel() string
}

type NotifyObject interface {
	GetCode() int
	SetCode(code int)
	GetMsg() string
	SetMsg(msg string)
	GetNotifyID() int
	GetSend() int
	SetSend(i int)
	GetReceive() int
	SetReceive(i int)
	GetDisplay() int
	SetDisplay(i int)
	GetClick() int
	SetClick(i int)
	GetValidDevice() int
	SetValidDevice(i int)
	GetActualSend() int
	SetActualSend(i int)
}

// PushService 提供推送服务的接口
type PushService interface {
	// Send 发送消息给单个设备
	Send(ctx context.Context, appid string, req interface{}, opt ...SendOption) error

	// GetNotifyStatus 查询通知发送状态
	GetNotifyStatus(ctx context.Context, appid string, notifyID string, obj NotifyObject) error

	// Name 获取推送的手机厂商名称
	Name() string
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	Success = 200
	Fail    = 400
)

type SendFunc func(ctx context.Context, token string) (*Response, error)

func RetrySend(ctx context.Context, send SendFunc, tokens []string, retry int, retryInterval int, maxConcurrent int) error {
	var wg sync.WaitGroup
	if retryInterval <= 0 {
		retryInterval = 1
	}
	if maxConcurrent <= 0 {
		maxConcurrent = 100
	}
	var MaxConcurrentPushes = make(chan struct{}, maxConcurrent)
	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentPushes <- struct{}{}
		wg.Add(1)
		go func(token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentPushes
				wg.Done()
			}()
			for i := 0; i <= retry; i++ {
				res, err := send(ctx, token)
				if err != nil || (res != nil && res.Code != 200) {
					if err == nil {
						err = errors.New(res.Msg)
					} else {
						es = append(es, err)
					}
					if i == 0 {
						continue
					}
					log.Printf("send error: %s (attempt %d)", err, i)
					time.Sleep(time.Duration(retryInterval) * time.Second)
				} else {
					log.Printf("send success: %s", res.Msg)
					break
				}
			}
		}(token)
	}
	wg.Wait()
	if len(es) > 0 {
		uniqueErrors := make(map[string]struct{})
		for _, err := range es {
			uniqueErrors[err.Error()] = struct{}{}
		}
		var uniqueErrorStrings []string
		for err := range uniqueErrors {
			uniqueErrorStrings = append(uniqueErrorStrings, err)
		}
		allErrorsString := strings.Join(uniqueErrorStrings, ", ")
		return errors.New(allErrorsString)
	}

	return nil
}
