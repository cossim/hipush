package config

type Config struct {
	Huawei HuaweiConfig `yaml:"huawei"`
}

type HuaweiConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppSecret string `yaml:"appsecret"`
	AppID     string `yaml:"appid"`
	Version   string `yaml:"version"`
	MaxRetry  int    `yaml:"max_retry"`
}

type XMConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppID     string `yaml:"appid"`
	AppSecret string `yaml:"appsecret"`
	MaxRetry  int    `yaml:"max_retry"`
}
