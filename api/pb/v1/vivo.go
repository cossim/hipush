package v1

func (x *VivoPushRequestData) GetCustomData() map[string]interface{} {
	return StructPBToMap(x.Data)
}

func (x *VivoPushRequestData) GetAppID() string {
	return x.Meta.AppID
}

func (x *VivoPushRequestData) GetAppName() string {
	return x.Meta.AppName
}

func (x *VivoPushRequestData) GetToken() []string {
	return x.Meta.Token
}

func (x *VivoPushRequestData) GetTopic() string {
	return ""
}

func (x *VivoPushRequestData) GetCollapseID() string {
	return ""
}

func (x *VivoPushRequestData) GetMessageID() string {
	return ""
}

func (x *VivoPushRequestData) GetPriority() string {
	return ""
}

func (x *VivoPushRequestData) GetCondition() string {
	return ""
}

func (x *VivoPushRequestData) GetIcon() string {
	return ""
}

func (x *VivoPushRequestData) GetMutableContent() bool {
	return false
}

func (x *VivoPushRequestData) GetContentAvailable() bool {
	return false
}
