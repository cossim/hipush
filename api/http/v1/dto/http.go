package dto

// PushRequest 表示推送请求的结构体
type PushRequest struct {
	// 推送平台 consts.Platform
	Platform string `json:"platform" binding:"required"`

	// 接收推送的设备标识
	Token []string `json:"token" binding:"required"`

	// ios capacitor.config文件中的appId 例如com.hitosea.apptest
	AppID string `json:"app_id" binding:"required"`

	AppSecret string `json:"app_secret"`

	// 自定义的消息数据，不同平台可能有不同的格式
	Data interface{} `json:"data"`
}
