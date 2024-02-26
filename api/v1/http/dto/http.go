package dto

// Platform 枚举表示推送平台
type Platform int

func (p Platform) String() string {
	switch p {
	case PlatformIOS:
		return "iOS"
	case PlatformHuawei:
		return "huawei"
	case PlatformGoogle:
		return "google"
	case PlatformXiaomi:
		return "xiaomi"
	case PlatformVivo:
		return "vivo"
	case PlatformOppo:
		return "oppo"
	default:
		return "Unknown"
	}
}

const (
	// PlatformIOS 表示 iOS 平台
	PlatformIOS Platform = iota + 1

	// PlatformHuawei 表示华为平台
	PlatformHuawei

	// PlatformGoogle 表示谷歌平台
	PlatformGoogle

	// PlatformXiaomi 表示小米平台
	PlatformXiaomi

	// PlatformVivo 表示 Vivo 平台
	PlatformVivo

	// PlatformOppo 表示 Oppo 平台
	PlatformOppo
)

// PushRequest 表示推送请求的结构体
type PushRequest struct {
	Platform  Platform    `json:"platform" binding:"required"` // 推送平台
	Token     []string    `json:"token" binding:"required"`    // 接收推送的设备标识
	AppID     string      `json:"app_id" binding:"required"`
	AppSecret string      `json:"app_secret" binding:"required"`
	Data      interface{} `json:"data"` // 自定义的消息数据，不同平台可能有不同的格式
}
