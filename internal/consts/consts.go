package consts

// Platform 枚举表示推送平台
type Platform int

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

// String converts the enum value to its string representation.
func (p Platform) String() string {
	switch p {
	case PlatformIOS:
		return "iOS"
	case PlatformHuawei:
		return "Huawei"
	case PlatformGoogle:
		return "Google"
	case PlatformXiaomi:
		return "Xiaomi"
	case PlatformVivo:
		return "Vivo"
	case PlatformOppo:
		return "Oppo"
	default:
		return "Unknown"
	}
}

// IsValid 判断平台是否有效
func (p Platform) IsValid() bool {
	switch p {
	case PlatformIOS, PlatformHuawei, PlatformGoogle, PlatformXiaomi, PlatformVivo, PlatformOppo:
		return true
	default:
		return false
	}
}
