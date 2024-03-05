package dto

// PushRequest 表示推送请求的结构体
type PushRequest struct {
	// AppID 应用程序标识
	// ios capacitor.config文件中的appId 例如com.hitosea.apptest
	AppID string `json:"app_id" binding:"required"`

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

// PushOption 表示推送选项的结构体
type PushOption struct {
	// DryRun 只进行数据校验不实际推送，数据校验成功即为成功
	DryRun bool `json:"dry_run,omitempty"`

	// Retry 重试次数
	Retry int `json:"retry,omitempty"`
}
