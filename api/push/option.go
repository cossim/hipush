package push

type SendOption interface {
	Apply(option *SendOptions)
}

// SendOptions 用于设置发送单个消息选项的结构体
type SendOptions struct {
	// DryRun 只进行数据校验不实际推送，数据校验成功即为成功
	DryRun bool `json:"dry_run,omitempty"`
	// Development 测试环境推送
	Development bool `json:"development,omitempty"`
	// Retry 重试次数
	Retry int32 `json:"retry,omitempty"`
	// RetryInterval 重试间隔（以秒为单位）
	RetryInterval int32 `json:"retry_interval"`
}

func (s *SendOptions) Apply(option *SendOptions) {
	if s.DryRun {
		option.DryRun = true
	}
	if s.Development {
		option.Development = true
	}
	option.Retry = s.Retry
	option.RetryInterval = s.RetryInterval
}

func (s *SendOptions) ApplyOptions(opts []SendOption) *SendOptions {
	for _, opt := range opts {
		opt.Apply(s)
	}
	return s
}

type MulticastOption interface {
	Apply(option *MulticastOptions)
}

// MulticastOptions 用于设置发送多个消息选项的结构体
type MulticastOptions struct {
	// MaxDevices 最大设备数
	MaxDevices int `json:"max_devices,omitempty"`
	// Delay 两次发送之间的延迟时间
	Delay int `json:"delay,omitempty"`
}

type SubscribeOption interface {
	Apply(option *SubscribeOptions)
}

// SubscribeOptions 用于设置订阅选项的结构体
type SubscribeOptions struct {
}

type UnsubscribeOption interface {
	Apply(option *UnsubscribeOption)
}

// UnsubscribeOptions 用于设置取消订阅选项的结构体
type UnsubscribeOptions struct {
}

type TopicOption interface {
	Apply(option *TopicOptions)
}

// TopicOptions 用于设置发送到特定主题选项的结构体
type TopicOptions struct {
}

type CheckDeviceOption interface {
	Apply(option *CheckDeviceOptions)
}

// CheckDeviceOptions 用于设置检查设备是否可用选项的结构体
type CheckDeviceOptions struct {
	// Timeout 检查设备的超时时间
	Timeout int `json:"timeout,omitempty"`
}
