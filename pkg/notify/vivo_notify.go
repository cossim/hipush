package notify

// VivoPushNotification
// https://dev.vivo.com.cn/documentCenter/doc/362#:~:text=%E6%8E%A5%E5%8F%A3%E5%AE%9A%E4%B9%89-,%E8%BE%93%E5%85%A5%E5%8F%82%E6%95%B0%EF%BC%9A,-intent%20uri
type VivoPushNotification struct {
	AppID     string `json:"app_id,omitempty"`
	RequestId string `json:"request_id,omitempty"`
	// Tokens 对应regId列表
	Tokens []string `json:"tokens" binding:"required"`

	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	Category string `json:"category,omitempty"`

	// Data 透传数据 客户端自定义键值对 key和Value键值对总长度不能超过1024字符
	Data map[string]string `json:"data,omitempty"`

	ClickAction *VivoClickAction `json:"click_action,omitempty"`

	// NotifyType 通知类型 1:无，2:响铃，3:振动，4:响铃和振动
	NotifyType int `json:"notify_type,omitempty"`

	// TTL 消息缓存时间，单位是秒，取值至少60秒，最长一天。当值为空时，默认一天
	TTL int `json:"ttl,omitempty"`

	// Retry 重试次数
	Retry int `json:"retry,omitempty"`

	// SendOnline true表示是在线直推，false表示非直推，设备离线直接丢弃
	SendOnline bool `json:"send_online,omitempty"`

	// Foreground 是否前台通知展示
	Foreground bool `json:"foreground,omitempty"`

	// Development 对应PushMode
	Development bool `json:"development,omitempty"`
}

type VivoClickAction struct {
	// Action 点击跳转类型 1：打开APP首页 2：打开链接 3：自定义 4:打开app内指定页面
	Action int `json:"action,omitempty"`

	// Activity 打开应用内页（activity 的 intent action）
	Activity string `json:"activity,omitempty"`

	// Url 打开网页的地址
	Url string `json:"url,omitempty"`
}
