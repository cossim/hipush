package v1

import "github.com/cossim/hipush/api/push"

var _ push.SendRequest = &AndroidPushRequestData{}

func (x *AndroidPushRequestData) GetNotifyType() int32 {
	return 0
}

func (x *AndroidPushRequestData) GetAppID() string {
	return x.Meta.AppID
}

func (x *AndroidPushRequestData) GetAppName() string {
	return x.Meta.AppName
}

func (x *AndroidPushRequestData) GetToken() []string {
	return x.Meta.Token
}

func (x *AndroidPushRequestData) GetMessageID() string {
	return ""
}

func (x *AndroidPushRequestData) GetCategory() string {
	return ""
}

func (x *AndroidPushRequestData) GetMutableContent() bool {
	return false
}

func (x *AndroidPushRequestData) GetContentAvailable() bool {
	return false
}

func (x *AndroidPushRequestData) GetDevelopment() bool {
	return false
}

func (x *AndroidPushRequestData) GetForeground() bool {
	return true
}
