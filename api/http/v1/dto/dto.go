package dto

import (
	"github.com/cossim/hipush/pkg/notify"
	"time"
)

type HuaweiPushRequestData struct {
	Title   string `json:"title"`
	Content string `json:"content"`

	// TTL represents the message cache time in seconds.
	// When the user device is offline, the message is cached on the Push server.
	// If the user device reconnects to the network within the message cache time, the message will be delivered.
	// After the cache time expires, the message will be discarded.
	// The default value is "86400s" (1 day), and the maximum value is "1296000s" (15 days).
	TTL string `json:"ttl,omitempty"`

	// Category The category of the notification
	// https://developer.huawei.com/consumer/cn/doc/HMSCore-References/https-send-api-0000001050986197#:~:text=%E8%BF%9B%E8%A1%8C%E7%BB%9F%E8%AE%A1%E5%88%86%E6%9E%90%E3%80%82-,category,-%E5%90%A6
	Category string `json:"category,omitempty"`

	// Priority The priority of the notification
	// normal、high default normal
	Priority string `json:"priority,omitempty"`

	// Icon Small icon URL
	Icon string `json:"icon,omitempty"`

	// Img Large image URL
	Img string `json:"img,omitempty"`

	// Sound represents the custom message notification ringtone.
	// It is effective when creating a new channel.
	// The ringtone file set here must be stored in the /res/raw path of the application.
	// For example, setting it to "/raw/shake" corresponds to the local "/res/raw/shake.xxx" file of the application.
	// Supported file formats include MP3, WAV, MPEG, etc.
	// If not set, the default system ringtone will be used.
	Sound string `json:"sound,omitempty"`

	// Foreground When the application is in the foreground, whether the notification bar message shows the switch
	Foreground bool `json:"foreground,omitempty"`

	// Development Test environment push
	Development bool `json:"development,omitempty"`

	// Type specifies the type of click action.
	// Possible values are:
	// 1: Open a custom app page.
	// 2: Open a specific URL.
	// 3: Open the app.
	// https://developer.huawei.com/consumer/cn/doc/HMSCore-References/https-send-api-0000001050986197#ZH-CN_TOPIC_0000001700731289__p431142991615:~:text=%E6%9C%80%E5%A4%A7%E9%95%BF%E5%BA%A61024-,ClickAction,-%E5%8F%82%E6%95%B0
	ClickAction ClickAction `json:"click_action,omitempty"`

	// https://developer.huawei.com/consumer/cn/doc/HMSCore-References/https-send-api-0000001050986197#ZH-CN_TOPIC_0000001700731289__p12819153131618:~:text=%E4%BA%8C%E9%80%89%E4%B8%80%E3%80%82-,BadgeNotification,-%E5%8F%82%E6%95%B0
	Badge BadgeNotification `json:"badge,omitempty"`
}

// BadgeNotification 结构体用于表示Android通知消息角标控制
type BadgeNotification struct {
	AddNum int    `json:"addNum,omitempty"`
	SetNum int    `json:"setNum,omitempty"`
	Class  string `json:"class"`
}

type APNsPushRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`

	// Topic The topic of the remote notification, which is typically the bundle ID
	// for your app. The certificate you create in the Apple Developer Member
	// Center must include the capability for this topic. If your certificate
	// includes multiple topics, you must specify a value for this header. If
	// you omit this header and your APNs certificate does not specify multiple
	// topics, the APNs server uses the certificate’s Subject as the default
	// topic.
	Topic string `json:"topic" binding:"required"`

	// CollapseID A string which allows multiple notifications with the same collapse
	// identifier to be displayed to the user as a single notification. The
	// value should not exceed 64 bytes.
	CollapseID string `json:"collapse_id,omitempty"`

	// An optional canonical UUID that identifies the notification. The
	// canonical form is 32 lowercase hexadecimal digits, displayed in five
	// groups separated by hyphens in the form 8-4-4-4-12. An example UUID is as
	// follows:
	//
	//  123e4567-e89b-12d3-a456-42665544000
	//
	// If you don't set this, a new UUID is created by APNs and returned in the
	// response.
	ApnsID string

	// Priority The priority of the notification
	// normal、high default normal
	Priority string `json:"priority,omitempty"`

	// PushType apns-push-type标头的值
	// https://developer.apple.com/documentation/usernotifications/sending-notification-requests-to-apns#Know-when-to-use-push-types
	PushType string `json:"push_type,omitempty"`

	URLArgs []string `json:"url_args,omitempty"`

	// TTL represents the expiration date of the notification.
	// If the value is nonzero, it indicates that the notification is valid until the specified date.
	// The value is a UNIX epoch expressed in seconds (UTC).
	// If the value is nonzero, APNs stores the notification and attempts to deliver it at least once, repeating the attempt as needed until the specified date.
	// If the value is 0, APNs attempts to deliver the notification only once and does not store it.
	TTL int64 `json:"TTL,omitempty"`

	// Badge sets the aps badge on the payload.
	// This will display a numeric badge on the app icon.
	Badge int `json:"badge,omitempty"`

	Development bool `json:"development,omitempty"`

	// MutableContent sets the aps mutable-content on the payload to 1.
	// This will indicate to the to the system to call your Notification Service
	// extension to mutate or replace the notification's content.
	MutableContent bool `json:"mutable_content,omitempty"`

	// ContentAvailable sets the aps content-available on the payload to 1.
	// This will indicate to the app that there is new content available to download
	// and launch the app in the background.
	ContentAvailable bool `json:"content_available,omitempty"`

	// Sound sets the aps sound on the payload.
	// This will play a sound from the app bundle, or the default sound otherwise.
	// https://developer.apple.com/documentation/usernotifications/generating-a-remote-notification Table 3. Keys to include in the sound dictionary
	Sound interface{} `json:"sound,omitempty"`

	// Data sets a custom key and value on the payload.
	// This will add custom key/value data to the notification payload at root level.
	Data map[string]interface{} `json:"data,omitempty"`
}

type Sound struct {
	Critical int     `json:"critical"`
	Name     string  `json:"name"`
	Volume   float64 `json:"volume"`
}

type VivoPushRequestData struct {
	Foreground bool `json:"foreground,omitempty"`
	// Development 测试环境推送
	Development bool `json:"development,omitempty"`
	TTL         int  `json:"ttl,omitempty"`
	// NotifyType 通知类型 1:无，2:响铃，3:振动，4:响铃和振动
	NotifyType  int               `json:"notify_type,omitempty"`
	NotifyID    int               `json:"notify_id,omitempty"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Category    string            `json:"category,omitempty"`
	TaskID      string            `json:"task_id,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	ClickAction ClickAction       `json:"click_action,omitempty"`
}

type OppoPushRequestData struct {
	Foreground bool   `json:"foreground,omitempty"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle,omitempty"`
	Content    string `json:"content"`
	// IsTimed 是否限时展示，指示消息是否在特定时间范围内展示
	IsTimed bool `json:"is_timed,omitempty"`
	// TimedDuration 限时展示时长，单位为秒，消息将在此时长内展示
	TimedDuration int `json:"timed_duration,omitempty"`
	// ValidityPeriod 消息有效时长，即推送服务缓存消息的时长，从消息创建是开始计算，最短为1小时，最长10天
	TTL int `json:"ttl,omitempty"`
	// IsScheduled false为立即推送 true为定时推送
	// 消息会在ScheduledStart-ScheduledEnd的时间段内随机展示
	IsScheduled bool `json:"is_scheduled,omitempty"`
	// ScheduledStart 定时推送的开始时间，指定消息推送的开始时间
	ScheduledStart time.Time `json:"scheduled_start"`
	// ScheduledEnd 定时推送的结束时间，指定消息推送的结束时间
	ScheduledEnd time.Time `json:"scheduled_end"`
	// Icon 消息图标，用于在通知栏上显示的图标
	Icon string `json:"icon,omitempty"`
	// ClickAction 点击动作
	ClickAction notify.OppoClickAction `json:"click_action"`
	// 附加的自定义参数
	Data map[string]string `json:"data,omitempty"`
}

type XiaomiPushRequestData struct {
	Title    string `json:"title,omitempty" binding:"required"`
	Subtitle string `json:"subtitle,omitempty"`
	Content  string `json:"content,omitempty" binding:"required"`

	// Foreground When the application is in the foreground, whether the notification bar message shows the switch
	Foreground bool `json:"foreground,omitempty"`

	// Icon 消息图标，用于在通知栏上显示的图标
	Icon string `json:"icon,omitempty"`

	// TTL 如果用户离线，设置消息在服务器保存的时间，单位：ms，服务器默认最长保留两周。
	TTL time.Duration `json:"ttl,omitempty"`

	// IsScheduled false为立即推送 true为定时推送
	// 消息会在ScheduledStart-ScheduledEnd的时间段内随机展示
	IsScheduled bool `json:"is_scheduled,omitempty"`
	// ScheduledTime 定时推送的开始时间，指定消息推送的开始时间
	// 用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间），仅支持七天内的定时消息。
	ScheduledTime time.Duration `json:"scheduled_time,omitempty"`

	// NotifyType represents the type of notification, and its value can be DEFAULT_ALL or a combination of the following:
	// DEFAULT_ALL = -1; DEFAULT_SOUND = 1;
	// Use the default sound for notification; DEFAULT_VIBRATE = 2;
	//Use default vibration for notification; DEFAULT_LIGHTS = 4;
	//Use default lights for notification.
	NotifyType int `json:"notify_type,omitempty"`

	// ClickAction Click behavior for predefined notification bar messages
	// "1": Open the Launcher Activity of the app after clicking on the notification in the notification bar.
	// "2": Open any Activity of the app after clicking on the notification in the notification bar (the developer also needs to pass url).
	// "3": Open a webpage after clicking on the notification in the notification bar.
	ClickAction ClickAction `json:"click_action"`

	// 附加的自定义参数
	Data map[string]string `json:"data,omitempty"`
}

type MeizuPushRequestData struct {
	Title   string `json:"title,omitempty" binding:"required"`
	Content string `json:"content,omitempty" binding:"required"`

	// TTL 如果用户离线，设置消息在服务器保存的时间，单位：ms，服务器默认最长保留两周。
	TTL int `json:"ttl,omitempty"`

	NotifyType int `json:"notify_type,omitempty" json:"notify_type,omitempty"`

	// Foreground 是否前台显示通知
	Foreground bool `json:"foreground,omitempty"`

	// IsScheduled 是否定时推送
	IsScheduled bool `json:"scheduled,omitempty"`
	// ScheduledStartTime 定时展示开始时间(yyyy-MM-dd HH:mm:ss)
	ScheduledStartTime string `json:"scheduled_start_time"`
	// ScheduledEndTime 定时展示结束时间(yyyy-MM-dd HH:mm:ss)
	ScheduledEndTime string `json:"scheduled_end_time"`

	// ClickAction 点击动作
	ClickAction ClickAction `json:"click_action"`

	// 附加的自定义参数
	Data map[string]string `json:"data,omitempty"`
}

type AndroidPushRequestData struct {
	Title   string `json:"title,omitempty" binding:"required"`
	Content string `json:"content,omitempty" binding:"required"`

	// TTL represents the duration for which the message is stored on the server if the user is offline.
	// The value should follow a specific format indicating the time duration, such as "86400s" for 1 day, "10m" for 10 minutes, or "1h" for 1 hour.
	TTL string `json:"ttl,omitempty"`

	Topic string `json:"topic,omitempty"`

	// Priority The priority of the notification
	// normal、high default normal
	Priority string `json:"priority,omitempty"`

	// CollapseID represents the collapse identifier of the notification.
	CollapseID string `json:"collapse_id,omitempty"`

	// Condition represents the condition for sending the notification to devices.
	Condition string `json:"condition,omitempty"`

	// Sound represents the custom sound for the push notification.
	Sound string `json:"sound,omitempty"`

	// Image represents the image associated with the push notification.
	Image string `json:"image,omitempty"`

	// 附加的自定义参数
	Data map[string]string `json:"data,omitempty"`
}

type HonorPushRequestData struct {
	Title   string `json:"title,omitempty" binding:"required"`
	Content string `json:"content,omitempty" binding:"required"`

	// Icon 消息图标，用于在通知栏上显示的图标
	Icon string `json:"icon,omitempty"`

	// Tag 消息标识，用于消息去重、覆盖
	Tag string `json:"tag,omitempty"`

	// Group 消息分组，例如发送10条带有同样group字段的消息，手机上只会展示该组消息中最新的一条和当前该组接收到的消息总数目，不会展示10条消息。
	Group string `json:"group,omitempty"`

	// NotifyId 消息通知ID，用于消息覆盖
	NotifyId int `json:"notify_id,omitempty"`

	// TTL 如果用户离线，设置消息在服务器保存的时间，单位：ms，服务器默认最长保留两周。
	TTL int `json:"ttl,omitempty"`

	NotifyType int `json:"notify_type,omitempty" json:"notify_type,omitempty"`

	// Development 测试模式推送消息
	Development bool `json:"development,omitempty"`

	// ClickAction 点击动作
	ClickAction ClickAction `json:"click_action"`

	// Badge 消息角标
	Badge BadgeNotification `json:"badge,omitempty"`

	// 附加的自定义参数
	Data map[string]interface{} `json:"data,omitempty"`
}

type ClickAction struct {
	// Action represents the click action.
	// Different manufacturers have different definitions.
	Action int `json:"action,omitempty"`

	// Activity opens an in-app page (activity's intent action).
	Activity string `json:"activity,omitempty"`

	// Url opens the URL of a webpage.
	Url string `json:"url,omitempty"`

	// Parameters represent the parameters appended to the URL after the URL redirection.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}
