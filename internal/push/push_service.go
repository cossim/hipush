package push

import "context"

type PushService interface {
	// 发送消息给单个设备
	Send(ctx context.Context, req interface{}) error

	// 发送消息给多个设备
	MulticastSend(ctx context.Context, req interface{}) error

	// 订阅特定主题
	Subscribe(ctx context.Context, req interface{}) error

	// 取消订阅特定主题
	Unsubscribe(ctx context.Context, req interface{}) error

	// 发送消息到特定主题
	SendToTopic(ctx context.Context, req interface{}) error

	// 发送消息到条件选择器
	SendToCondition(ctx context.Context, req interface{}) error

	// 检查设备是否可用
	CheckDevice(ctx context.Context, req interface{}) bool

	// 获取推送服务的名称
	Name() string
}
