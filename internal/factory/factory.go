package factory

import (
	"errors"
	"github.com/cossim/hipush/push"
)

type PushServiceCreator func() push.PushService

type PushServiceFactory struct {
	creators map[string]push.PushService
}

func NewPushServiceFactory() *PushServiceFactory {
	return &PushServiceFactory{
		creators: make(map[string]push.PushService),
	}
}

func (f *PushServiceFactory) Register(name string, creator PushServiceCreator) {
	f.creators[name] = creator()
}

func (f *PushServiceFactory) GetPushService(name string) (push.PushService, error) {
	creator, ok := f.creators[name]
	if !ok {
		return nil, errors.New("unsupported platform")
	}
	return creator, nil
}
