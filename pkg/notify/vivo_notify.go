package notify

// VivoPushNotification
// https://dev.vivo.com.cn/documentCenter/doc/362#:~:text=%E6%8E%A5%E5%8F%A3%E5%AE%9A%E4%B9%89-,%E8%BE%93%E5%85%A5%E5%8F%82%E6%95%B0%EF%BC%9A,-intent%20uri
type VivoPushNotification struct {
	AppID     string `json:"app_id,omitempty"`
	RequestId string `json:"request_id,omitempty"`

	// NotifyID 每条消息在通知栏的唯一标识 可以用于覆盖消息
	NotifyID int `json:"notify_id,omitempty"`

	TaskID string `json:"task_id,omitempty"`

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

// VivoPushStats 表示 Vivo 推送统计信息
type VivoPushStats struct {
	Code        int    `json:"code"`         // 接口调用状态码
	Msg         string `json:"msg"`          // 文字描述接口调用情况
	NotifyID    int    `json:"notify_id"`    // 消息ID
	Send        int    `json:"send"`         // 发送量
	Receive     int    `json:"receive"`      // 到达量
	Display     int    `json:"display"`      // 展示量
	Click       int    `json:"click"`        // 点击量
	ValidDevice int    `json:"valid_device"` // 有效设备量
	ActualSend  int    `json:"actual_send"`  // 实际发送量
}

// GetCode 返回接口调用状态码
func (v *VivoPushStats) GetCode() int {
	return v.Code
}

// SetCode 设置接口调用状态码
func (v *VivoPushStats) SetCode(code int) {
	v.Code = code
}

// GetMsg 返回文字描述接口调用情况
func (v *VivoPushStats) GetMsg() string {
	return v.Msg
}

// SetMsg 设置文字描述接口调用情况
func (v *VivoPushStats) SetMsg(msg string) {
	v.Msg = msg
}

// GetNotifyID 返回消息ID
func (v *VivoPushStats) GetNotifyID() int {
	return v.NotifyID
}

// GetSend 返回发送量
func (v *VivoPushStats) GetSend() int {
	return v.Send
}

// SetSend 设置发送量
func (v *VivoPushStats) SetSend(i int) {
	v.Send = i
}

// GetReceive 返回到达量
func (v *VivoPushStats) GetReceive() int {
	return v.Receive
}

// SetReceive 设置到达量
func (v *VivoPushStats) SetReceive(i int) {
	v.Receive = i
}

// GetDisplay 返回展示量
func (v *VivoPushStats) GetDisplay() int {
	return v.Display
}

// SetDisplay 设置展示量
func (v *VivoPushStats) SetDisplay(i int) {
	v.Display = i
}

// GetClick 返回点击量
func (v *VivoPushStats) GetClick() int {
	return v.Click
}

// SetClick 设置点击量
func (v *VivoPushStats) SetClick(i int) {
	v.Click = i
}

// GetValidDevice 返回有效设备量
func (v *VivoPushStats) GetValidDevice() int {
	return v.ValidDevice
}

// SetValidDevice 设置有效设备量
func (v *VivoPushStats) SetValidDevice(i int) {
	v.ValidDevice = i
}

// GetActualSend 返回实际发送量
func (v *VivoPushStats) GetActualSend() int {
	return v.ActualSend
}

// SetActualSend 设置实际发送量
func (v *VivoPushStats) SetActualSend(i int) {
	v.ActualSend = i
}
