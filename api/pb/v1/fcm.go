package v1

func (x *AndroidPushRequestData) GetCustomData() map[string]interface{} {
	return StructPBToMap(x.Data)
}

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
