package adapter

import (
	"context"
	"github.com/cossim/hipush/push"
)

func NewPushServiceAdapter(pushService push.PushService) *PushServiceAdapter {
	return &PushServiceAdapter{pushService: pushService}
}

type PushServiceAdapter struct {
	pushService push.PushService
}

func (p *PushServiceAdapter) Send(ctx context.Context, req interface{}, opt ...push.SendOption) error {
	return p.pushService.Send(ctx, req, opt...)
}

func (p *PushServiceAdapter) SendMulticast(ctx context.Context, req interface{}, opt ...push.MulticastOption) error {
	//TODO implement me
	panic("implement me")
}

func (p *PushServiceAdapter) Subscribe(ctx context.Context, req interface{}, opt ...push.SubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (p *PushServiceAdapter) Unsubscribe(ctx context.Context, req interface{}, opt ...push.UnsubscribeOption) error {
	//TODO implement me
	panic("implement me")
}

func (p *PushServiceAdapter) SendToTopic(ctx context.Context, req interface{}, opt ...push.TopicOption) error {
	//TODO implement me
	panic("implement me")
}

func (p *PushServiceAdapter) CheckDevice(ctx context.Context, req interface{}, opt ...push.CheckDeviceOption) bool {
	//TODO implement me
	panic("implement me")
}

func (p *PushServiceAdapter) Name() string {
	//TODO implement me
	panic("implement me")
}
