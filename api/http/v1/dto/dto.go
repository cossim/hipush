package dto

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
