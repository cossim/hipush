package notify

type HMSPushNotification struct {
	AppID         string      `json:"app_id,omitempty"`
	AppName       string      `json:"app_name,omitempty"`
	Tokens        []string    `json:"tokens" binding:"required"`
	Topic         string      `json:"topic,omitempty"`
	Title         string      `json:"title,omitempty"`
	Category      string      `json:"category,omitempty"`
	Content       string      `json:"content,omitempty"`
	Condition     string      `json:"condition,omitempty"`
	Priority      string      `json:"priority,omitempty"`
	CollapseKey   int         `json:"collapse_key,omitempty"`
	Image         string      `json:"image,omitempty"`
	Data          string      `json:"huawei_data,omitempty"`
	TTL           string      `json:"huawei_ttl,omitempty"`
	BiTag         string      `json:"bi_tag,omitempty"`
	Sound         interface{} `json:"sound,omitempty"`
	FastAppTarget int         `json:"fast_app_target,omitempty"`

	// Badge 消息角标
	Badge BadgeNotification `json:"badge,omitempty"`

	Development bool `json:"development,omitempty"`
}
