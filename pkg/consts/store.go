package consts

const (
	key = "hipush"

	// 总数
	TotalSuffix = "-total"
	// 成功数
	SuccessSuffix = "-success"
	// 失败数
	FailedSuffix = "-Failed"

	SendSuffix    = "-send"
	ReceiveSuffix = "-receive"
	DisplaySuffix = "-display"
	ClickSuffix   = "-click"
)

// 为每个平台生成键名称
const (
	// HTTP方式
	HTTPPrefix = key + "-http"
	// GRPC方式
	GRPCPrefix = key + "-pb"

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
	// 荣耀平台
	HonorPrefix = key + "-honor"
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
	IosSend    = iOSPrefix + SendSuffix
	IosReceive = iOSPrefix + ReceiveSuffix
	IosDisplay = iOSPrefix + DisplaySuffix
	IosClick   = iOSPrefix + ClickSuffix

	// 华为平台键名称
	HuaweiTotal   = HuaweiPrefix + TotalSuffix
	HuaweiSuccess = HuaweiPrefix + SuccessSuffix
	HuaweiFailed  = HuaweiPrefix + FailedSuffix
	HuaweiSend    = HuaweiPrefix + SendSuffix
	HuaweiReceive = HuaweiPrefix + ReceiveSuffix
	HuaweiDisplay = HuaweiPrefix + DisplaySuffix
	HuaweiClick   = HuaweiPrefix + ClickSuffix

	// 安卓平台键名称
	AndroidTotal   = AndroidPrefix + TotalSuffix
	AndroidSuccess = AndroidPrefix + SuccessSuffix
	AndroidFailed  = AndroidPrefix + FailedSuffix
	AndroidSend    = AndroidPrefix + SendSuffix
	AndroidReceive = AndroidPrefix + ReceiveSuffix
	AndroidDisplay = AndroidPrefix + DisplaySuffix
	AndroidClick   = AndroidPrefix + ClickSuffix

	// Vivo平台键名称
	VivoTotal   = VivoPrefix + TotalSuffix
	VivoSuccess = VivoPrefix + SuccessSuffix
	VivoFailed  = VivoPrefix + FailedSuffix
	VivoSend    = VivoPrefix + SendSuffix
	VivoReceive = VivoPrefix + ReceiveSuffix
	VivoDisplay = VivoPrefix + DisplaySuffix
	VivoClick   = VivoPrefix + ClickSuffix

	// Oppo平台键名称
	OppoTotal   = OppoPrefix + TotalSuffix
	OppoSuccess = OppoPrefix + SuccessSuffix
	OppoFailed  = OppoPrefix + FailedSuffix
	OppoSend    = OppoPrefix + SendSuffix
	OppoReceive = OppoPrefix + ReceiveSuffix
	OppoDisplay = OppoPrefix + DisplaySuffix
	OppoClick   = OppoPrefix + ClickSuffix

	// 小米平台键名称
	XiaomiTotal   = XiaomiPrefix + TotalSuffix
	XiaomiSuccess = XiaomiPrefix + SuccessSuffix
	XiaomiFailed  = XiaomiPrefix + FailedSuffix
	XiaomiSend    = XiaomiPrefix + SendSuffix
	XiaomiReceive = XiaomiPrefix + ReceiveSuffix
	XiaomiDisplay = XiaomiPrefix + DisplaySuffix
	XiaomiClick   = XiaomiPrefix + ClickSuffix

	// 魅族平台键名称
	MeizuTotal   = MeizuPrefix + TotalSuffix
	MeizuSuccess = MeizuPrefix + SuccessSuffix
	MeizuFailed  = MeizuPrefix + FailedSuffix
	MeizuSend    = MeizuPrefix + SendSuffix
	MeizuReceive = MeizuPrefix + ReceiveSuffix
	MeizuDisplay = MeizuPrefix + DisplaySuffix
	MeizuClick   = MeizuPrefix + ClickSuffix

	// 荣耀平台键名称
	HonorTotal   = HonorPrefix + TotalSuffix
	HonorSuccess = HonorPrefix + SuccessSuffix
	HonorFailed  = HonorPrefix + FailedSuffix
	HonorSend    = HonorPrefix + SendSuffix
	HonorReceive = HonorPrefix + ReceiveSuffix
	HonorDisplay = HonorPrefix + DisplaySuffix
	HonorClick   = HonorPrefix + ClickSuffix
)

// iOS平台
var ()

// 华为平台
var (
	huaweiSend    = HuaweiPrefix + "-send"
	huaweiReceive = HuaweiPrefix + "-receive"
	huaweiDisplay = HuaweiPrefix + "-display"
	huaweiClick   = HuaweiPrefix + "-click"
)

// 安卓平台
var (
	androidSend    = AndroidPrefix + "-send"
	androidReceive = AndroidPrefix + "-receive"
	androidDisplay = AndroidPrefix + "-display"
	androidClick   = AndroidPrefix + "-click"
)

// Vivo平台
var ()

// Oppo平台
var (
	oppoSend    = OppoPrefix + "-send"
	oppoReceive = OppoPrefix + "-receive"
	oppoDisplay = OppoPrefix + "-display"
	oppoClick   = OppoPrefix + "-click"
)

// 小米平台
var (
	xiaomiSend    = XiaomiPrefix + "-send"
	xiaomiReceive = XiaomiPrefix + "-receive"
	xiaomiDisplay = XiaomiPrefix + "-display"
	xiaomiClick   = XiaomiPrefix + "-click"
)

// 魅族平台
var (
	meizuSend    = MeizuPrefix + "-send"
	meizuReceive = MeizuPrefix + "-receive"
	meizuDisplay = MeizuPrefix + "-display"
	meizuClick   = MeizuPrefix + "-click"
)

// 荣耀平台
var (
	honorSend    = HonorPrefix + "-send"
	honorReceive = HonorPrefix + "-receive"
	honorDisplay = HonorPrefix + "-display"
	honorClick   = HonorPrefix + "-click"
)
