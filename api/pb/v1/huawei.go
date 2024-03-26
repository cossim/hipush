package v1

import "github.com/cossim/hipush/api/push"

var _ push.SendRequest = &HuaweiPushRequestData{}

func (x *HuaweiPushRequestData) GetNotifyType() int32 {
	return 0
}

func (x *HuaweiPushRequestData) GetAppID() string {
	return x.Meta.AppID
}

func (x *HuaweiPushRequestData) GetAppName() string {
	return x.Meta.AppName
}

func (x *HuaweiPushRequestData) GetToken() []string {
	return x.Meta.Token
}

func (x *HuaweiPushRequestData) GetTopic() string {
	return ""
}

func (x *HuaweiPushRequestData) GetCollapseID() string {
	return ""
}

func (x *HuaweiPushRequestData) GetMessageID() string {
	return ""
}

func (x *HuaweiPushRequestData) GetCondition() string {
	return ""
}

func (x *HuaweiPushRequestData) GetMutableContent() bool {
	return false
}

func (x *HuaweiPushRequestData) GetContentAvailable() bool {
	return false
}
