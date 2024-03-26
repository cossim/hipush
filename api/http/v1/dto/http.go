package dto

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
	GRPC    PushStat `json:"pb"`      // GRPC 推送状态
}
