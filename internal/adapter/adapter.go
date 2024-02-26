package adapter

import (
	"context"
	"github.com/cossim/hipush/internal/push"
)

type PushServiceAdapter struct {
	pushService push.PushService
}

func NewPushServiceAdapter(pushService push.PushService) *PushServiceAdapter {
	return &PushServiceAdapter{pushService: pushService}
}

func (p *PushServiceAdapter) Send(ctx context.Context, req interface{}) error {
	return p.pushService.Send(ctx, req)
}

func (p *PushServiceAdapter) MulticastSend(ctx context.Context, req interface{}) error {
	return p.pushService.MulticastSend(ctx, req)
}

func (p *PushServiceAdapter) Subscribe(ctx context.Context, req interface{}) error {
	return p.pushService.Subscribe(ctx, req)
}

func (p *PushServiceAdapter) Unsubscribe(ctx context.Context, req interface{}) error {
	return p.pushService.Unsubscribe(ctx, req)
}

func (p *PushServiceAdapter) SendToTopic(ctx context.Context, req interface{}) error {
	return p.pushService.SendToTopic(ctx, req)
}

func (p *PushServiceAdapter) SendToCondition(ctx context.Context, req interface{}) error {
	return p.pushService.SendToCondition(ctx, req)
}

func (p *PushServiceAdapter) CheckDevice(ctx context.Context, req interface{}) bool {
	return p.pushService.CheckDevice(ctx, req)
}

func (p *PushServiceAdapter) Name() string {
	return p.pushService.Name()
}
