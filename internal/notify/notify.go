package notify

import (
	"github.com/appleboy/go-fcm"
	"github.com/cossim/go-hms-push/push/model"
)

const (
	HIGH   = "high"
	NORMAL = "nornal"
)

type ClickAction struct {
	Action  string `json:"action,omitempty"`
	Content string `json:"content,omitempty"`
}

// PushOption 表示推送选项的结构体
type PushOption struct {
	// DryRun 只进行数据校验不实际推送，数据校验成功即为成功
	DryRun bool `json:"dry_run,omitempty"`

	// Retry 重试次数
	Retry int `json:"retry,omitempty"`
}

// D provide string array
type D map[string]interface{}

// Sound sets the aps sound on the payload.
// https://developer.apple.com/documentation/usernotifications/generating-a-remote-notification#:~:text=%E8%AD%A6%E6%8A%A5%E7%9A%84%E5%A3%B0%E9%9F%B3%E3%80%82-,sound%E8%A1%A8%203.%E5%AD%97%E5%85%B8%E4%B8%AD%E5%8C%85%E5%90%AB%E7%9A%84%E9%94%AE,-%E9%92%A5%E5%8C%99
type Sound struct {
	// Critical 指示声音是否被标记为关键声音。关键声音通常用于需要立即引起用户注意的通知。
	// 值为 1 表示是关键声音，值为 0 表示不是关键声音。默认为 0。
	Critical int `json:"critical,omitempty"`
	// Name 声音的名称或标识符。
	// 通常是声音文件的名称，表示要播放的声音文件。默认为空字符串。
	Name string `json:"name,omitempty"`
	// Volume 声音的音量级别。
	// 值范围为 0.0 到 1.0，表示音量的相对级别。默认为 1.0。
	Volume float32 `json:"volume,omitempty"`
}

// PushNotification is single notification request
type PushNotification struct {
	// Common
	ID               string      `json:"notif_id,omitempty"`
	Tokens           []string    `json:"tokens" binding:"required"`
	Platform         int         `json:"platform" binding:"required"`
	Message          string      `json:"message,omitempty"`
	Title            string      `json:"title,omitempty"`
	Image            string      `json:"image,omitempty"`
	Priority         string      `json:"priority,omitempty"`
	ContentAvailable bool        `json:"content_available,omitempty"`
	MutableContent   bool        `json:"mutable_content,omitempty"`
	Sound            interface{} `json:"sound,omitempty"`
	Data             D           `json:"data,omitempty"`
	Retry            int         `json:"retry,omitempty"`

	// Android
	APIKey                string            `json:"api_key,omitempty"`
	To                    string            `json:"to,omitempty"`
	CollapseKey           string            `json:"collapse_key,omitempty"`
	DelayWhileIdle        bool              `json:"delay_while_idle,omitempty"`
	TimeToLive            *uint             `json:"time_to_live,omitempty"`
	RestrictedPackageName string            `json:"restricted_package_name,omitempty"`
	DryRun                bool              `json:"dry_run,omitempty"`
	Condition             string            `json:"condition,omitempty"`
	Notification          *fcm.Notification `json:"notification,omitempty"`

	// Huawei
	AppID              string                     `json:"app_id,omitempty"`
	AppSecret          string                     `json:"app_secret,omitempty"`
	HuaweiNotification *model.AndroidNotification `json:"huawei_notification,omitempty"`
	HuaweiData         string                     `json:"huawei_data,omitempty"`
	HuaweiCollapseKey  int                        `json:"huawei_collapse_key,omitempty"`
	HuaweiTTL          string                     `json:"huawei_ttl,omitempty"`
	BiTag              string                     `json:"bi_tag,omitempty"`
	FastAppTarget      int                        `json:"fast_app_target,omitempty"`

	// iOS
	Expiration  *int64   `json:"expiration,omitempty"`
	ApnsID      string   `json:"apns_id,omitempty"`
	CollapseID  string   `json:"collapse_id,omitempty"`
	Topic       string   `json:"topic,omitempty"`
	PushType    string   `json:"push_type,omitempty"`
	Badge       *int     `json:"badge,omitempty"`
	Category    string   `json:"category,omitempty"`
	ThreadID    string   `json:"thread-id,omitempty"`
	URLArgs     []string `json:"url-args,omitempty"`
	Alert       Alert    `json:"alert,omitempty"`
	Production  bool     `json:"production,omitempty"`
	Development bool     `json:"development,omitempty"`
	SoundName   string   `json:"name,omitempty"`
	SoundVolume float32  `json:"volume,omitempty"`
	Apns        D        `json:"apns,omitempty"`

	// ref: https://github.com/sideshow/apns2/blob/54928d6193dfe300b6b88dad72b7e2ae138d4f0a/payload/builder.go#L7-L24
	InterruptionLevel string `json:"interruption_level,omitempty"`
}
