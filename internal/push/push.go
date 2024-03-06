package push

import (
	"context"
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

type SendOption struct {
	// DryRun 只进行数据校验不实际推送，数据校验成功即为成功
	DryRun bool `json:"dry_run,omitempty"`
	// Retry 重试次数
	Retry int `json:"retry,omitempty"`
}

type PushService interface {
	// 发送消息给单个设备
	//Send(ctx context.Context, req PushRequest) error

	Send(ctx context.Context, req interface{}, opt SendOption) error

	// 发送消息给多个设备
	MulticastSend(ctx context.Context, req interface{}) error

	// 订阅特定主题
	Subscribe(ctx context.Context, req interface{}) error

	// 取消订阅特定主题
	Unsubscribe(ctx context.Context, req interface{}) error

	// 发送消息到特定主题
	SendToTopic(ctx context.Context, req interface{}) error

	// 发送消息到条件选择器
	SendToCondition(ctx context.Context, req interface{}) error

	// 检查设备是否可用
	CheckDevice(ctx context.Context, req interface{}) bool

	// 获取推送服务的名称
	Name() string
}
