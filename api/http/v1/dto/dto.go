package dto

import (
	"github.com/cossim/hipush/internal/notify"
	"time"
)

// ClickAction 点击行为
type ClickAction struct {
	Url    string
	Action Action // 点击行为类型
}

// Action 枚举表示点击行为类型
type Action int

const (
	// ActionTypeOpenCustomPage 表示打开应用自定义页面
	ActionTypeOpenCustomPage Action = iota + 1

	// ActionTypeOpenURL 表示点击后打开特定URL
	ActionTypeOpenURL

	// ActionTypeOpenApp 表示点击后打开应用
	ActionTypeOpenApp
)

type HuaweiPushRequestData struct {
	DryRun      bool
	Foreground  bool
	TTL         string
	Type        string
	Title       string
	Message     string
	Category    string
	Icon        string
	Img         string
	Sound       string
	ClickAction ClickAction // 点击行为
}

type APNsPushRequest struct {
	Topic            string                 `json:"topic,omitempty"`
	CollapseID       string                 `json:"collapse_id,omitempty"`
	Priority         string                 `json:"priority,omitempty"`
	PushType         string                 `json:"push_type,omitempty"`
	Title            string                 `json:"title,omitempty"`
	Message          string                 `json:"message,omitempty"`
	SoundName        string                 `json:"sound_name,omitempty"`
	ThreadID         string                 `json:"thread_id,omitempty"`
	URLArgs          []string               `json:"url_args,omitempty"`
	Expiration       *int64                 `json:"expiration,omitempty"`
	Badge            *int                   `json:"badge,omitempty"`
	SoundVolume      float64                `json:"sound_volume,omitempty"`
	Production       bool                   `json:"production,omitempty"`
	Development      bool                   `json:"development,omitempty"`
	MutableContent   bool                   `json:"mutable_content,omitempty"`
	ContentAvailable bool                   `json:"content_available,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
	// 声音配置（可选）
	// https://developer.apple.com/documentation/usernotifications/generating-a-remote-notification Table 3. Keys to include in the sound dictionary
	Sound interface{} `json:"sound,omitempty"`
}

type Sound struct {
	Critical int     `json:"critical"`
	Name     string  `json:"name"`
	Volume   float64 `json:"volume"`
}

type VivoPushRequestData struct {
	DryRun     bool
	Foreground bool
	TTL        int
	Type       string
	Title      string
	Message    string
	Category   string
	Data       map[string]string
}

type OppoPushRequestData struct {
	Foreground bool
	Title      string
	Subtitle   string
	Message    string
	// IsTimed 是否限时展示，指示消息是否在特定时间范围内展示
	IsTimed bool
	// TimedDuration 限时展示时长，单位为秒，消息将在此时长内展示
	TimedDuration int
	// ValidityPeriod 消息有效时长，即推送服务缓存消息的时长，从消息创建是开始计算，最短为1小时，最长10天
	ValidityPeriod int
	// IsScheduled false为立即推送 true为定时推送
	// 消息会在ScheduledStart-ScheduledEnd的时间段内随机展示
	IsScheduled bool
	// ScheduledStart 定时推送的开始时间，指定消息推送的开始时间
	ScheduledStart time.Time
	// ScheduledEnd 定时推送的结束时间，指定消息推送的结束时间
	ScheduledEnd time.Time
	// Icon 消息图标，用于在通知栏上显示的图标
	Icon string
	// ClickAction 点击动作
	ClickAction notify.OppoClickAction
	// 附加的自定义参数
	Data map[string]string
}
