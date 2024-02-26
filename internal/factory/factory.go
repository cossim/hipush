package factory

import (
	"errors"
	"github.com/cossim/hipush/internal/push"
)

type PushServiceCreator func() push.PushService

type PushServiceFactory struct {
	creators map[string]PushServiceCreator
}

func NewPushServiceFactory() *PushServiceFactory {
	return &PushServiceFactory{
		creators: make(map[string]PushServiceCreator),
	}
}

func (f *PushServiceFactory) Register(name string, creator PushServiceCreator) {
	f.creators[name] = creator
}

func (f *PushServiceFactory) CreatePushService(name string) (push.PushService, error) {
	creator, ok := f.creators[name]
	if !ok {
		return nil, errors.New("unsupported platform")
	}
	return creator(), nil
}
