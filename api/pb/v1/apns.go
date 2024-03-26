package v1

import (
	"github.com/cossim/hipush/api/push"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"time"
)

var _ push.SendRequest = &APNsPushRequest{}

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

func (m *APNsPushRequest) BuildNotification(req push.SendRequest) (*apns2.Notification, error) {
	topic := req.GetTopic()
	if topic == "" {
		topic = req.GetAppID()
	}
	notification := &apns2.Notification{
		ApnsID:     req.GetMessageID(),
		Topic:      topic,
		CollapseID: req.GetCollapseID(),
	}

	if req.GetTTL() != 0 {
		notification.Expiration = time.Unix(req.GetTTL(), 0)
	}

	if req.GetPriority() == "normal" {
		notification.Priority = apns2.PriorityLow
	} else if req.GetPriority() == "high" {
		notification.Priority = apns2.PriorityHigh
	}

	//if len(req.PushType) > 0 {
	//	notification.PushType = apns2.EPushType(req.PushType)
	//}

	payload := payload.NewPayload()

	// add alert object if message length > 0 and title is empty
	if len(req.GetContent()) > 0 && req.GetTitle() == "" {
		payload.Alert(req.GetContent())
	}

	// zero value for clear the badge on the app icon.
	//if req.Badge != nil && *req.Badge >= 0 {
	//	payload.Badge(*req.Badge)
	//}

	if req.GetMutableContent() {
		payload.MutableContent()
	}

	//switch req.GetSound().(type) {
	//// from http request binding
	//case map[string]interface{}:
	//	result := &Sound{}
	//	_ = mapstructure.Decode(req.GetSound(), &result)
	//	payload.Sound(result)
	//// from http request binding for non critical alerts
	//case string:
	//	payload.Sound(req.GetSound())
	//case Sound:
	//	payload.Sound(req.GetSound())
	//}

	//if len(req.SoundName) > 0 {
	//	payload.SoundName(req.SoundName)
	//}
	//
	//if req.SoundVolume > 0 {
	//	payload.SoundVolume(req.SoundVolume)
	//}

	if req.GetContentAvailable() {
		payload.ContentAvailable()
	}

	//if len(req.URLArgs) > 0 {
	//	payload.URLArgs(req.URLArgs)
	//}

	//if len(req.ThreadID) > 0 {
	//	payload.ThreadID(req.ThreadID)
	//}

	//for k, v := range req.GetData() {
	//	payload.Custom(k, v)
	//}

	payload.AlertTitle(req.GetTitle())
	payload.AlertBody(req.GetContent())
	payload.Category(req.GetCategory())

	//payload = iosAlertDictionary(payload, req)

	notification.Payload = payload

	return notification, nil
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
