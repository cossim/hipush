package dto

// PushRequest 表示推送请求的结构体
type PushRequest struct {
	// AppID 应用程序标识
	// ios capacitor.config文件中的appId 例如com.hitosea.apptest
	AppID string `json:"app_id"`

	// appName 应用名称
	AppName string `json:"app_name"`

	// Platform 推送平台 consts.Platform
	Platform string `json:"platform" binding:"required"`

	// Token 接收推送的设备标识
	// 例如ios为deviceToken
	// vivo、oppo为RegId
	Token []string `json:"token" binding:"required"`

	// 推送的消息请求数据，不同平台可能有不同的格式
	Data interface{} `json:"data" binding:"required"`

	// PushOptions 推送选项
	Option PushOption `json:"option,omitempty"`
}

type PushOption struct {
	// DryRun 只进行数据校验不实际推送，数据校验成功即为成功
	DryRun bool `json:"dry_run,omitempty"`

	// Retry 重试次数
	Retry int `json:"retry,omitempty"`

	// RetryInterval 重试间隔（以秒为单位）
	RetryInterval int `json:"retry_interval,omitempty"`
}

type PushMessageStatRequest struct {
	// Platform 平台名称 consts.Platform
	Platform string `json:"platform" binding:"required"`

	AppName string `json:"app_name"`

	AppID string `json:"app_id"`

	TaskID []string `json:"task_id" binding:"required"`
}

// PushStat 每个推送平台的推送状态
type PushStat struct {
	Total   int64 `json:"total"`   // 总推送数
	Success int64 `json:"success"` // 成功推送数
	Failed  int64 `json:"failed"`  // 失败推送数
	Send    int64 `json:"send"`    // 发送数
	Receive int64 `json:"receive"` // 到达数
	Display int64 `json:"display"` // 展示数
	Click   int64 `json:"click"`   // 点击数
}

// PushStats 所有推送平台的推送状态
type PushStats struct {
	// PushStat 所有平台总推送数据
	PushStat

	Android PushStat `json:"android"` // Android 平台推送状态
	IOS     PushStat `json:"ios"`     // iOS 平台推送状态
	Xiaomi  PushStat `json:"xiaomi"`  // 小米平台推送状态
	Vivo    PushStat `json:"vivo"`    // Vivo 平台推送状态
	Oppo    PushStat `json:"oppo"`    // Oppo 平台推送状态
	Meizu   PushStat `json:"meizu"`   // 魅族平台推送状态
	Huawei  PushStat `json:"huawei"`  // 华为平台推送状态
	Honor   PushStat `json:"honor"`   // 荣耀平台推送状态
	HTTP    PushStat `json:"http"`    // HTTP 推送状态
	GRPC    PushStat `json:"grpc"`    // GRPC 推送状态
}
