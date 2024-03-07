package consts

const (
	key = "hipush"

	// 总数
	TotalSuffix = "-total"
	// 成功数
	SuccessSuffix = "-success"
	// 失败数
	FailedSuffix = "-Failed"
)

// 为每个平台生成键名称
const (
	// HTTP方式
	HTTPPrefix = key + "-http"
	// GRPC方式
	GRPCPrefix = key + "-grpc"

	// iOS平台
	iOSPrefix = key + "-ios"
	// 华为平台
	HuaweiPrefix = key + "-huawei"
	// 安卓平台
	AndroidPrefix = key + "-android"
	// Vivo平台
	VivoPrefix = key + "-vivo"
	// Oppo平台
	OppoPrefix = key + "-oppo"
	// 小米平台
	XiaomiPrefix = key + "-xiaomi"
	// 魅族平台
	MeizuPrefix = key + "-meizu"
)

var (
	HiPushTotal   = key + "-total"
	HiPushSuccess = key + SuccessSuffix
	HiPushFailed  = key + SuccessSuffix

	HTTPTotal   = HTTPPrefix + TotalSuffix
	HTTPSuccess = HTTPPrefix + SuccessSuffix
	HTTPFailed  = HTTPPrefix + FailedSuffix

	GRPCTotal   = GRPCPrefix + TotalSuffix
	GRPCSuccess = GRPCPrefix + SuccessSuffix
	GRPCFailed  = GRPCPrefix + FailedSuffix

	// iOS平台键名称
	IosTotal   = iOSPrefix + TotalSuffix
	IosSuccess = iOSPrefix + SuccessSuffix
	IosFailed  = iOSPrefix + FailedSuffix

	// 华为平台键名称
	HuaweiTotal   = HuaweiPrefix + TotalSuffix
	HuaweiSuccess = HuaweiPrefix + SuccessSuffix
	HuaweiFailed  = HuaweiPrefix + FailedSuffix

	// 安卓平台键名称
	AndroidTotal   = AndroidPrefix + TotalSuffix
	AndroidSuccess = AndroidPrefix + SuccessSuffix
	AndroidFailed  = AndroidPrefix + FailedSuffix

	// Vivo平台键名称
	VivoTotal   = VivoPrefix + TotalSuffix
	VivoSuccess = VivoPrefix + SuccessSuffix
	VivoFailed  = VivoPrefix + FailedSuffix

	// Oppo平台键名称
	OppoTotal   = OppoPrefix + TotalSuffix
	OppoSuccess = OppoPrefix + SuccessSuffix
	OppoFailed  = OppoPrefix + FailedSuffix

	// 小米平台键名称
	XiaomiTotal   = XiaomiPrefix + TotalSuffix
	XiaomiSuccess = XiaomiPrefix + SuccessSuffix
	XiaomiFailed  = XiaomiPrefix + FailedSuffix

	// 魅族平台键名称
	MeizuTotal   = MeizuPrefix + TotalSuffix
	MeizuSuccess = MeizuPrefix + SuccessSuffix
	MeizuFailed  = MeizuPrefix + FailedSuffix
)
