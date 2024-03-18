package notify

type FCMPushNotification struct {
	AppID            string   `json:"app_id,omitempty"`
	AppName          string   `json:"app_name,omitempty"`
	Tokens           []string `json:"tokens" binding:"required"`
	Topic            string   `json:"topic,omitempty"`
	Priority         string   `json:"priority,omitempty"`
	Title            string   `json:"title,omitempty"`
	Message          string   `json:"message,omitempty"`
	Image            string   `json:"image,omitempty"`
	Sound            string   `json:"sound,omitempty"`
	CollapseID       string   `json:"collapse_id,omitempty"`
	Category         string   `json:"category,omitempty"`
	Condition        string   `json:"condition,omitempty"`
	TTL              *uint    `json:"ttl,omitempty"`
	Retry            int      `json:"retry,omitempty"`
	Badge            *int     `json:"badge,omitempty"`
	ContentAvailable bool     `json:"content_available,omitempty"`
	MutableContent   bool     `json:"mutable_content"`
	Data             D        `json:"data,omitempty"`
	Apns             D        `json:"apns,omitempty"`
}
