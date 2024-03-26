package v1

func (m *APNsPushRequest) GetCustomData() map[string]interface{} {
	return StructPBToMap(m.Data)
}

func (m *APNsPushRequest) GetNotifyType() int32 {
	return 0
}

func (m *APNsPushRequest) GetAppID() string {
	return m.Meta.AppID
}

func (m *APNsPushRequest) GetAppName() string {
	return m.Meta.AppName
}

func (m *APNsPushRequest) GetToken() []string {
	return m.Meta.Token
}

func (m *APNsPushRequest) GetMessageID() string {
	return m.ApnsID
}

func (m *APNsPushRequest) GetCondition() string {
	return ""
}

func (m *APNsPushRequest) GetIcon() string {
	return ""
}

func (m *APNsPushRequest) GetForeground() bool {
	return true
}

// Sound sets the aps sound on the payload.
// https://developer.apple.com/documentation/usernotifications/generating-a-remote-notification#:~:text=%E8%AD%A6%E6%8A%A5%E7%9A%84%E5%A3%B0%E9%9F%B3%E3%80%82-,sound%E8%A1%A8%203.%E5%AD%97%E5%85%B8%E4%B8%AD%E5%8C%85%E5%90%AB%E7%9A%84%E9%94%AE,-%E9%92%A5%E5%8C%99
type Sound struct {
	// Critical 指示声音是否被标记为关键声音。关键声音通常用于需要立即引起用户注意的通知。
	// 值为 1 表示是关键声音，值为 0 表示不是关键声音。默认为 0。
	Critical int `json:"critical,omitempty"`
	// Name 声音的名称或标识符。
	// 通常是声音文件的名称，表示要播放的声音文件。默认为空字符串。
	Name string `json:"name,omitempty"`
	// Volume 声音的音量级别。
	// 值范围为 0.0 到 1.0，表示音量的相对级别。默认为 1.0。
	Volume float32 `json:"volume,omitempty"`
}
