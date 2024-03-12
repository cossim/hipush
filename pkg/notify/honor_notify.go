package notify

type HonorPushNotification struct {
	AppID    string   `json:"app_id,omitempty"`
	Tokens   []string `json:"tokens" binding:"required"`
	Title    string   `json:"title,omitempty"`
	Content  string   `json:"content,omitempty"`
	Image    string   `json:"image,omitempty"`
	Priority string   `json:"priority,omitempty"`
	Category string   `json:"category,omitempty"`
	TTL      string   `json:"ttl,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`

	Development bool `json:"development,omitempty"`

	Badge *BadgeNotification `json:"badge,omitempty"`

	ClickAction *HonorClickAction `json:"click_action,omitempty"`
}

// BadgeNotification 结构体用于表示Android通知消息角标控制
type BadgeNotification struct {
	AddNum     int    `json:"addNum,omitempty"`
	SetNum     int    `json:"setNum,omitempty"`
	BadgeClass string `json:"badgeClass"`
}

type HonorClickAction struct {
	// Action 点击跳转类型
	// 1 打开应用内页（activity的action标签名）
	// 2 打开特定url
	// 3 打开应用
	Action int `json:"action,omitempty"`

	// Activity 打开应用内页（activity 的 intent action）
	Activity string `json:"activity,omitempty"`

	// Url 打开网页的地址
	Url string `json:"url,omitempty"`

	// Parameters url跳转后传的参数拼接在url后面
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}
