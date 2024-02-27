package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type iOSAppConfig struct {
	MaxConcurrentPushes int    `yaml:"max_concurrent_pushes"`
	MaxRetry            int    `yaml:"max_retry"`
	Enabled             bool   `yaml:"enabled"`
	Production          bool   `yaml:"production"`
	AppID               string `yaml:"appid"`
	KeyPath             string `yaml:"key_path"`
	KeyType             string `yaml:"key_type"`
	Password            string `yaml:"password"`
	KeyID               string `yaml:"key_id"`
	TeamID              string `yaml:"team_id"`
}

type HuaweiAppConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppID     string `yaml:"appid"`
	AppSecret string `yaml:"appsecret"`
	MaxRetry  int    `yaml:"max_retry"`
}

type Config struct {
	IOS    []iOSAppConfig    `yaml:"ios"`
	Huawei []HuaweiAppConfig `yaml:"huawei"`
}

func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
