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

func (f *PushServiceFactory) WithPushService(ps push.PushService) PushServiceCreator {
	return func() push.PushService {
		return ps
	}
}

func (f *PushServiceFactory) Register(creators ...PushServiceCreator) error {
	for _, c := range creators {
		ps := c()
		f.creators[ps.Name()] = ps
	}
	return nil
}

func (f *PushServiceFactory) GetPushService(name string) (push.PushService, error) {
	ps, ok := f.creators[name]
	if !ok {
		return nil, errors.New("unsupported platform")
	}
	return ps, nil
}
