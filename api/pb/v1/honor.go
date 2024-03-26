package v1

func (x *HonorPushRequestData) GetCustomData() map[string]interface{} {
	return StructPBToMap(x.Data)
}

func (x *HonorPushRequestData) GetAppID() string {
	return x.Meta.AppID
}

func (x *HonorPushRequestData) GetAppName() string {
	return x.Meta.AppName
}

func (x *HonorPushRequestData) GetToken() []string {
	return x.Meta.Token
}

func (x *HonorPushRequestData) GetTopic() string {
	return ""
}

func (x *HonorPushRequestData) GetCollapseID() string {
	return ""
}

func (x *HonorPushRequestData) GetMessageID() string {
	return ""
}

func (x *HonorPushRequestData) GetPriority() string {
	return ""
}

func (x *HonorPushRequestData) GetCategory() string {
	return ""
}

func (x *HonorPushRequestData) GetCondition() string {
	return ""
}

func (x *HonorPushRequestData) GetMutableContent() bool {
	return false
}

func (x *HonorPushRequestData) GetContentAvailable() bool {
	return false
}

func (x *HonorPushRequestData) GetForeground() bool {
	return true
}

func (x *HonorPushRequestData) GetNotifyType() int32 {
	return 0
}
