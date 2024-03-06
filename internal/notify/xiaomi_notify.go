package notify

import "time"

// XiaomiPushNotification
// https://dev.mi.com/console/doc/detail?pId=2776#_0
type XiaomiPushNotification struct {
	AppID string `json:"app_id,omitempty"`

	// Tokens 对应regId列表
	Tokens []string `json:"tokens" binding:"required"`

	Title   string `json:"title,omitempty"`
	Content string `json:"message,omitempty"`

	// NotifyType
	// DEFAULT_ALL = -1
	// DEFAULT_SOUND = 1; 使用默认提示音提示
	// DEFAULT_VIBRATE = 2; 使用默认振动提示
	// DEFAULT_LIGHTS = 4; 使用默认呼吸灯提示。
	NotifyType int `json:"notify_type,omitempty"`

	// TTL 如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
	TTL int64 `json:"time_to_live,omitempty"`

	IsShowNotify bool `json:"is_show_notify,omitempty"`

	IsScheduled bool
	// ScheduledTime 定时推送的开始时间，指定消息推送的开始时间
	// 用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间），仅支持七天内的定时消息。
	ScheduledTime time.Duration
}
