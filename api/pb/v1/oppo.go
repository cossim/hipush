package v1

func (x *OppoPushRequestData) GetCustomData() map[string]interface{} {
	return StructPBToMap(x.Data)
}

func (x *OppoPushRequestData) GetAppID() string {
	return x.Meta.AppID
}

func (x *OppoPushRequestData) GetAppName() string {
	return x.Meta.AppName
}

func (x *OppoPushRequestData) GetToken() []string {
	return x.Meta.Token
}

func (x *OppoPushRequestData) GetTopic() string {
	return ""
}

func (x *OppoPushRequestData) GetCollapseID() string {
	return ""
}

func (x *OppoPushRequestData) GetMessageID() string {
	return ""
}

func (x *OppoPushRequestData) GetPriority() string {
	return ""
}

func (x *OppoPushRequestData) GetCategory() string {
	return ""
}

func (x *OppoPushRequestData) GetCondition() string {
	return ""
}

func (x *OppoPushRequestData) GetMutableContent() bool {
	return false
}

func (x *OppoPushRequestData) GetContentAvailable() bool {
	return false
}
