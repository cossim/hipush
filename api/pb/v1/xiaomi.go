package v1

func (x *XiaomiPushRequestData) GetCustomData() map[string]interface{} {
	return StructPBToMap(x.Data)
}

func (x *XiaomiPushRequestData) GetAppID() string {
	return x.Meta.AppID
}

func (x *XiaomiPushRequestData) GetAppName() string {
	return x.Meta.AppName
}

func (x *XiaomiPushRequestData) GetToken() []string {
	return x.Meta.Token
}

func (x *XiaomiPushRequestData) GetTopic() string {
	return ""
}

func (x *XiaomiPushRequestData) GetCollapseID() string {
	return ""
}

func (x *XiaomiPushRequestData) GetMessageID() string {
	return ""
}

func (x *XiaomiPushRequestData) GetPriority() string {
	return ""
}

func (x *XiaomiPushRequestData) GetCategory() string {
	return ""
}

func (x *XiaomiPushRequestData) GetCondition() string {
	return ""
}

func (x *XiaomiPushRequestData) GetMutableContent() bool {
	return false
}

func (x *XiaomiPushRequestData) GetContentAvailable() bool {
	return false
}
