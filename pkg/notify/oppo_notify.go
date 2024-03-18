package notify

type OppoPushNotification struct {
	AppID     string `json:"app_id"`
	AppName   string `json:"app_name"`
	RequestId string `json:"request_id,omitempty"`
	// Tokens 对应regId列表
	Tokens []string `json:"tokens" binding:"required"`

	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Message  string `json:"message,omitempty"`

	// Data 透传数据 客户端自定义键值对 key和Value键值对总长度不能超过1024字符
	Data map[string]string `json:"data,omitempty"`

	ClickAction *OppoClickAction `json:"click_action,omitempty"`

	// TTL 消息缓存时间，单位是秒，取值至少60秒，最长一天。当值为空时，默认一天
	TTL int `json:"ttl,omitempty"`

	// PushOptions 推送选项
	Option PushOption `json:"option,omitempty"`
}

type OppoClickAction struct {
	// Action 点击跳转类型 1：打开APP首页 2：打开链接 3：自定义 4:打开app内指定页面 5:跳转Intentscheme URL   默认值为 0
	// 0 启动应用
	// 1 打开应用内页（activity的action标签名）
	// 2 打开网页
	// 4 打开应用内页（activity 全路径类名）
	// 5 Intentscheme URL
	Action int `json:"action,omitempty"`

	// Activity 打开应用内页（activity 的 intent action）
	Activity string `json:"activity,omitempty"`

	// Url 打开网页的地址
	Url string `json:"url,omitempty"`

	// Parameters url跳转后传的参数拼接在url后面
	Parameters string `json:"parameters,omitempty"`
	Content    string `json:"content,omitempty"`
}
