package factory

import (
	"errors"
	"github.com/cossim/hipush/push"
)

type PushServiceCreator func() (push.PushService, error)

type PushServiceFactory struct {
	creators map[string]push.PushService
}

func NewPushServiceFactory() *PushServiceFactory {
	return &PushServiceFactory{
		creators: make(map[string]push.PushService),
	}
}

func WithPushServiceCreator(creator func() (push.PushService, error)) PushServiceCreator {
	return func() (push.PushService, error) {
		ps, err := creator()
		if err != nil {
			return nil, err
		}
		return ps, nil
	}
}

func (f *PushServiceFactory) Register(creators ...PushServiceCreator) error {
	for _, c := range creators {
		ps, err := c()
		if err != nil {
			return err
		}
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
