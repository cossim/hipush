package notify

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"time"
)

type ApnsPushNotification struct {
	AppID string `json:"app_id,omitempty"`

	Tokens            []string               `json:"tokens" binding:"required"`
	Priority          string                 `json:"priority,omitempty"`
	Title             string                 `json:"title,omitempty"`
	Content           string                 `json:"content,omitempty"`
	Expiration        *int64                 `json:"expiration,omitempty"`
	ApnsID            string                 `json:"apns_id,omitempty"`
	CollapseID        string                 `json:"collapse_id,omitempty"`
	Topic             string                 `json:"topic,omitempty"`
	PushType          string                 `json:"push_type,omitempty"`
	Badge             *int                   `json:"badge,omitempty"`
	Category          string                 `json:"category,omitempty"`
	ThreadID          string                 `json:"thread-id,omitempty"`
	URLArgs           []string               `json:"url-args,omitempty"`
	Alert             Alert                  `json:"alert,omitempty"`
	ContentAvailable  bool                   `json:"content_available,omitempty"`
	MutableContent    bool                   `json:"mutable_content"`
	Production        bool                   `json:"production,omitempty"`
	Development       bool                   `json:"development,omitempty"`
	SoundName         string                 `json:"name,omitempty"`
	SoundVolume       float32                `json:"volume,omitempty"`
	Apns              D                      `json:"apns,omitempty"`
	InterruptionLevel string                 `json:"interruption_level,omitempty"`
	Sound             interface{}            `json:"sound,omitempty"`
	Data              map[string]interface{} `json:"data,omitempty"`
}

func (a *ApnsPushNotification) Get() interface{} {
	//TODO implement me
	panic("implement me")
}

func (a *ApnsPushNotification) GetRetry() int {
	return 0
}

func (a *ApnsPushNotification) GetTokens() []string {
	return a.Tokens
}

func (a *ApnsPushNotification) GetTitle() string {
	return a.Title
}

func (a *ApnsPushNotification) GetMessage() string {
	return a.Content
}

func (a *ApnsPushNotification) GetTopic() string {
	return a.Topic
}

func (a *ApnsPushNotification) GetKey() string {
	//TODO implement me
	panic("implement me")
}

func (a *ApnsPushNotification) GetCategory() string {
	return a.Category
}

func (a *ApnsPushNotification) GetSound() interface{} {
	return a.Sound
}

func (a *ApnsPushNotification) GetAlert() interface{} {
	return a.Alert
}

func (a *ApnsPushNotification) GetBadge() int {
	if a.Badge == nil {
		return 0
	}
	return *a.Badge
}

func (a *ApnsPushNotification) GetThreadID() string {
	return a.ThreadID
}

func (a *ApnsPushNotification) GetData() map[string]interface{} {
	return a.Data
}

func (a *ApnsPushNotification) GetImage() string {
	//TODO implement me
	panic("implement me")
}

func (a *ApnsPushNotification) GetID() string {
	//TODO implement me
	panic("implement me")
}

func (a *ApnsPushNotification) GetPushType() string {
	return a.PushType
}

func (a *ApnsPushNotification) GetPriority() string {
	return a.Priority
}

func (a *ApnsPushNotification) IsContentAvailable() bool {
	return a.ContentAvailable
}

func (a *ApnsPushNotification) IsMutableContent() bool {
	return a.MutableContent
}

func (a *ApnsPushNotification) IsDevelopment() bool {
	return a.Development
}

func (a *ApnsPushNotification) GetExpiration() *int64 {
	if a.Expiration == nil {
		return nil
	}
	return a.Expiration
}

func (a *ApnsPushNotification) GetApnsID() string {
	return a.ApnsID
}

func (a *ApnsPushNotification) GetCollapseID() string {
	return a.CollapseID
}

func (a *ApnsPushNotification) GetURLArgs() []string {
	return a.URLArgs
}

func (a *ApnsPushNotification) GetSoundName() string {
	return a.SoundName
}

func (a *ApnsPushNotification) GetSoundVolume() float32 {
	return a.SoundVolume
}

func (a *ApnsPushNotification) GetApns() map[string]interface{} {
	return a.Apns
}

func (a *ApnsPushNotification) GetInterruptionLevel() string {
	return a.InterruptionLevel
}

// Alert is APNs payload
type Alert struct {
	Action          string   `json:"action,omitempty"`
	ActionLocKey    string   `json:"action-loc-key,omitempty"`
	Body            string   `json:"body,omitempty"`
	LaunchImage     string   `json:"launch-image,omitempty"`
	LocArgs         []string `json:"loc-args,omitempty"`
	LocKey          string   `json:"loc-key,omitempty"`
	Title           string   `json:"title,omitempty"`
	Subtitle        string   `json:"subtitle,omitempty"`
	TitleLocArgs    []string `json:"title-loc-args,omitempty"`
	TitleLocKey     string   `json:"title-loc-key,omitempty"`
	SummaryArg      string   `json:"summary-arg,omitempty"`
	SummaryArgCount int      `json:"summary-arg-count,omitempty"`
}

// GetIOSNotification use for define iOS notification.
// The iOS Notification Payload (Payload Key Reference)
// Ref: https://apple.co/2VtH6Iu
func GetIOSNotification(req *ApnsPushNotification) *apns2.Notification {
	notification := &apns2.Notification{
		//ApnsID:     req.ApnsID,
		Topic:      req.Topic,
		CollapseID: req.CollapseID,
	}

	if req.Expiration != nil {
		notification.Expiration = time.Unix(*req.Expiration, 0)
	}

	if len(req.Priority) > 0 {
		if req.Priority == "normal" {
			notification.Priority = apns2.PriorityLow
		} else if req.Priority == HIGH {
			notification.Priority = apns2.PriorityHigh
		}
	}

	if len(req.PushType) > 0 {
		notification.PushType = apns2.EPushType(req.PushType)
	}

	payload := payload.NewPayload()

	// add alert object if message length > 0 and title is empty
	if len(req.Content) > 0 && req.Title == "" {
		payload.Alert(req.Content)
	}

	// zero value for clear the badge on the app icon.
	if req.Badge != nil && *req.Badge >= 0 {
		payload.Badge(*req.Badge)
	}

	if req.MutableContent {
		payload.MutableContent()
	}

	switch req.Sound.(type) {
	// from http request binding
	case map[string]interface{}:
		result := &Sound{}
		_ = mapstructure.Decode(req.Sound, &result)
		payload.Sound(result)
	// from http request binding for non critical alerts
	case string:
		payload.Sound(&req.Sound)
	case Sound:
		payload.Sound(&req.Sound)
	}

	if len(req.SoundName) > 0 {
		payload.SoundName(req.SoundName)
	}

	if req.SoundVolume > 0 {
		payload.SoundVolume(req.SoundVolume)
	}

	if req.ContentAvailable {
		payload.ContentAvailable()
	}

	if len(req.URLArgs) > 0 {
		payload.URLArgs(req.URLArgs)
	}

	if len(req.ThreadID) > 0 {
		payload.ThreadID(req.ThreadID)
	}

	for k, v := range req.Data {
		payload.Custom(k, v)
	}

	payload = iosAlertDictionary(payload, req)

	notification.Payload = payload

	return notification
}

func iosAlertDictionary(notificationPayload *payload.Payload, req *ApnsPushNotification) *payload.Payload {
	// Alert dictionary

	if len(req.Title) > 0 {
		notificationPayload.AlertTitle(req.Title)
	}

	if len(req.InterruptionLevel) > 0 {
		notificationPayload.InterruptionLevel(payload.EInterruptionLevel(req.InterruptionLevel))
	}

	if len(req.Content) > 0 && len(req.Title) > 0 {
		notificationPayload.AlertBody(req.Content)
	}

	if len(req.Alert.Title) > 0 {
		notificationPayload.AlertTitle(req.Alert.Title)
	}

	// Apple Watch & Safari display this string as part of the notification interface.
	if len(req.Alert.Subtitle) > 0 {
		notificationPayload.AlertSubtitle(req.Alert.Subtitle)
	}

	if len(req.Alert.TitleLocKey) > 0 {
		notificationPayload.AlertTitleLocKey(req.Alert.TitleLocKey)
	}

	if len(req.Alert.LocArgs) > 0 {
		notificationPayload.AlertLocArgs(req.Alert.LocArgs)
	}

	if len(req.Alert.TitleLocArgs) > 0 {
		notificationPayload.AlertTitleLocArgs(req.Alert.TitleLocArgs)
	}

	if len(req.Alert.Body) > 0 {
		notificationPayload.AlertBody(req.Alert.Body)
	}

	if len(req.Alert.LaunchImage) > 0 {
		notificationPayload.AlertLaunchImage(req.Alert.LaunchImage)
	}

	if len(req.Alert.LocKey) > 0 {
		notificationPayload.AlertLocKey(req.Alert.LocKey)
	}

	if len(req.Alert.Action) > 0 {
		notificationPayload.AlertAction(req.Alert.Action)
	}

	if len(req.Alert.ActionLocKey) > 0 {
		notificationPayload.AlertActionLocKey(req.Alert.ActionLocKey)
	}

	// General
	if len(req.Category) > 0 {
		notificationPayload.Category(req.Category)
	}

	if len(req.Alert.SummaryArg) > 0 {
		notificationPayload.AlertSummaryArg(req.Alert.SummaryArg)
	}
	if req.Alert.SummaryArgCount > 0 {
		notificationPayload.AlertSummaryArgCount(req.Alert.SummaryArgCount)
	}

	return notificationPayload
}
