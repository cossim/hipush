package consts

// Platform 枚举表示推送平台
type Platform string

const (
	// PlatformUnknown 表示未知平台
	PlatformUnknown Platform = "unknown"

	// PlatformIOS 表示 iOS 平台
	PlatformIOS Platform = "ios"

	// PlatformAndroid 表示安卓平台
	PlatformAndroid Platform = "android"

	// PlatformHuawei 表示华为平台
	PlatformHuawei Platform = "huawei"

	// PlatformXiaomi 表示小米平台
	PlatformXiaomi Platform = "xiaomi"

	// PlatformVivo 表示 Vivo 平台
	PlatformVivo Platform = "vivo"

	// PlatformOppo 表示 Oppo 平台
	PlatformOppo Platform = "oppo"

	// PlatformMeizu 表示魅族平台
	PlatformMeizu Platform = "meizu"

	// PlatformHonor 表示荣耀平台
	PlatformHonor Platform = "honor"
)

// PlatformSlice 存储所有平台的切片
var PlatformSlice = []Platform{
	PlatformIOS,
	PlatformAndroid,
	PlatformHuawei,
	PlatformXiaomi,
	PlatformVivo,
	PlatformOppo,
	PlatformMeizu,
	PlatformHonor,
}

// String converts the enum value to its string representation.
func (p Platform) String() string {
	switch p {
	case PlatformIOS:
		return "ios"
	case PlatformAndroid:
		return "android"
	case PlatformHuawei:
		return "huawei"
	case PlatformXiaomi:
		return "xiaomi"
	case PlatformVivo:
		return "vivo"
	case PlatformOppo:
		return "oppo"
	case PlatformMeizu:
		return "meizu"
	default:
		return "unknown"
	}
}

// IsValid 判断平台是否有效
func (p Platform) IsValid() bool {
	switch p {
	case PlatformIOS, PlatformHuawei, PlatformAndroid, PlatformXiaomi, PlatformVivo, PlatformOppo, PlatformMeizu, PlatformHonor:
		return true
	default:
		return false
	}
}
