package notify

import "time"

type FCMPushNotification struct {
	AppID      string            `json:"app_id,omitempty"`
	AppName    string            `json:"app_name,omitempty"`
	Tokens     []string          `json:"tokens,omitempty"`
	Title      string            `json:"title,omitempty"`
	Content    string            `json:"content,omitempty"`
	Topic      string            `json:"topic,omitempty"`
	Priority   string            `json:"priority,omitempty"`
	Image      string            `json:"image,omitempty"`
	Sound      string            `json:"sound,omitempty"`
	CollapseID string            `json:"collapse_id,omitempty"`
	Category   string            `json:"category,omitempty"`
	Condition  string            `json:"condition,omitempty"`
	TTL        time.Duration     `json:"ttl,omitempty"`
	Badge      *int              `json:"badge,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
	Apns       D                 `json:"apns,omitempty"`
}
