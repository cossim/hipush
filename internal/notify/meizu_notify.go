package notify

// MeizuPushNotification
// https://github.com/MEIZUPUSH/PushAPI/blob/master/README.md
type MeizuPushNotification struct {
	AppID string `json:"app_id,omitempty"`

	// Tokens 对应pushId列表
	Tokens []string `json:"tokens" binding:"required"`

	Title   string `json:"title,omitempty"`
	Content string `json:"message,omitempty"`

	// NotifyType
	// DEFAULT_ALL = -1;
	// DEFAULT_SOUND = 0; 使用默认提示音提示
	// DEFAULT_VIBRATE = 1; 使用默认振动提示
	// DEFAULT_LIGHTS = 2; 使用默认呼吸灯提示。
	NotifyType int `json:"notify_type,omitempty"`

	ClickAction *MeizuClickAction `json:"click_action,omitempty"`

	// TTL 如果用户离线，设置消息在服务器保存的时间，有效时长 (1 -72 小时内的正整数)
	TTL int `json:"time_to_live,omitempty"`

	OffLine bool `json:"offline,omitempty"`

	IsShowNotify bool `json:"is_show_notify,omitempty"`

	// IsScheduled 是否定时推送
	IsScheduled bool
	// ScheduledStartTime 定时展示开始时间(yyyy-MM-dd HH:mm:ss)
	ScheduledStartTime string
	// ScheduledEndTime 定时展示结束时间(yyyy-MM-dd HH:mm:ss)
	ScheduledEndTime string
}

type MeizuClickAction struct {
	// Action 点击跳转类型
	// 0 打开应用
	// 1 打开应用内页（activity的action标签名）
	// 2 打开H5地址（应用本地的URI）
	Action int `json:"action,omitempty"`

	// Activity 打开应用内页（activity 的 intent action）
	// 格式 pkg.activity eg: com.meizu.upspushdemo.TestActivity
	Activity string `json:"activity,omitempty"`

	// Url 打开网页的地址
	Url string `json:"url,omitempty"`

	// Parameters url跳转后传的参数拼接在url后面
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}
