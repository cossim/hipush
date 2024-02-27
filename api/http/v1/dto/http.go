package dto

// Platform 枚举表示推送平台
type Platform int

// PushRequest 表示推送请求的结构体
type PushRequest struct {
	Platform Platform `json:"platform" binding:"required"` // 推送平台
	Token    []string `json:"token" binding:"required"`    // 接收推送的设备标识
	// ios capacitor.config文件中的appId 例如com.hitosea.apptest
	AppID     string      `json:"app_id" binding:"required"`
	AppSecret string      `json:"app_secret" binding:"required"`
	Data      interface{} `json:"data"` // 自定义的消息数据，不同平台可能有不同的格式
}
