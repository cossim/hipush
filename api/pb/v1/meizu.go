package v1

import "github.com/cossim/hipush/api/push"

var _ push.SendRequest = &MeizuPushRequestData{}

func (x *MeizuPushRequestData) GetAppID() string {
	return x.Meta.AppID
}

func (x *MeizuPushRequestData) GetAppName() string {
	return x.Meta.AppName
}

func (x *MeizuPushRequestData) GetToken() []string {
	return x.Meta.Token
}

func (x *MeizuPushRequestData) GetTopic() string {
	return ""
}

func (x *MeizuPushRequestData) GetCollapseID() string {
	return ""
}

func (x *MeizuPushRequestData) GetMessageID() string {
	return ""
}

func (x *MeizuPushRequestData) GetPriority() string {
	return ""
}

func (x *MeizuPushRequestData) GetCategory() string {
	return ""
}

func (x *MeizuPushRequestData) GetCondition() string {
	return ""
}

func (x *MeizuPushRequestData) GetIcon() string {
	return ""
}

func (x *MeizuPushRequestData) GetMutableContent() bool {
	return false
}

func (x *MeizuPushRequestData) GetContentAvailable() bool {
	return false
}
