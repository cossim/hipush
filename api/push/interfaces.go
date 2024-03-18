package push

import (
	"context"
)

type TaskObjectList interface {
	Add(obj TaskObject)
	Get() []TaskObject
}

type TaskObject interface {
	GetCode() int
	SetCode(code int)
	GetMsg() string
	SetMsg(msg string)
	GetTaskID() string
	SetTaskID(id string)
	GetSend() int
	SetSend(i int)
	GetReceive() int
	SetReceive(i int)
	GetDisplay() int
	SetDisplay(i int)
	GetClick() int
	SetClick(i int)
	GetInvalidDevice() int
	SetInvalidDevice(i int)
	GetValidDevice() int
	SetValidDevice(i int)
}

type SendResponse struct {
	TaskId string `json:"task_id"`
}

// PushService 提供推送服务的接口
type PushService interface {
	// Send 发送消息给单个设备
	Send(ctx context.Context, req interface{}, opt ...SendOption) (*SendResponse, error)

	// GetTasksStatus 查询推送消息的统计信息
	GetTasksStatus(ctx context.Context, appid string, taskID []string, obj TaskObjectList) error

	// Name 获取推送的手机厂商名称
	Name() string
}
