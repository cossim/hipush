package push

type PushMessageStatsList struct {
	Items []TaskObject `json:"items"`
}

func (v *PushMessageStatsList) Add(obj TaskObject) {
	v.Items = append(v.Items, obj)
}

func (v *PushMessageStatsList) Get() []TaskObject {
	return v.Items
}

// VivoPushStats 表示 Vivo 推送统计信息
type VivoPushStats struct {
	TaskID        string `json:"task_id"`        // 消息的任务id
	Msg           string `json:"msg"`            // 文字描述接口调用情况
	Code          int    `json:"code"`           // 接口调用状态码
	Send          int    `json:"send"`           // 发送量
	Receive       int    `json:"receive"`        // 到达量
	Display       int    `json:"display"`        // 展示量
	Click         int    `json:"click"`          // 点击量
	InvalidDevice int    `json:"invalid_device"` // 无效设备量
	ValidDevice   int    `json:"valid_device"`   // 有效设备量
}

func (v *VivoPushStats) GetInvalidDevice() int {
	return v.InvalidDevice
}

func (v *VivoPushStats) SetInvalidDevice(i int) {
	v.InvalidDevice = i
}

// XiaomiPushStats 表示 Vivo 推送统计信息
type XiaomiPushStats struct {
	TaskID      string `json:"task_id"`      // 消息的任务id
	Msg         string `json:"msg"`          // 文字描述接口调用情况
	Code        int    `json:"code"`         // 接口调用状态码
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
func (v *VivoPushStats) GetTaskID() string {
	return v.TaskID
}

func (v *VivoPushStats) SetTaskID(id string) {
	v.TaskID = id
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
