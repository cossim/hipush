package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type iOSAppConfig struct {
	MaxConcurrentPushes int    `yaml:"max_concurrent_pushes"`
	MaxRetry            int    `yaml:"max_retry"`
	Enabled             bool   `yaml:"enabled"`
	Production          bool   `yaml:"production"`
	AppID               string `yaml:"app_id"`
	KeyPath             string `yaml:"key_path"`
	KeyType             string `yaml:"key_type"`
	Password            string `yaml:"password"`
	KeyID               string `yaml:"key_id"`
	TeamID              string `yaml:"team_id"`
}

type HuaweiAppConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
	AuthUrl   string `yaml:"auth_url"`
	PushUrl   string `yaml:"push_url"`
	MaxRetry  int    `yaml:"max_retry"`
}

type VivoAppConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppID     string `yaml:"app_id"`
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
	MaxRetry  int    `yaml:"max_retry"`
}

type OppoAppConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppID     string `yaml:"app_id"`
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
	MaxRetry  int    `yaml:"max_retry"`
}

type AndroidAppConfig struct {
	Enabled  bool   `yaml:"enabled"`
	AppID    string `yaml:"app_id"`
	AppKey   string `yaml:"app_key"`
	MaxRetry int    `yaml:"max_retry"`
}

type XiaomiAppConfig struct {
	Enabled   bool     `yaml:"enabled"`
	AppID     string   `yaml:"app_id"`
	Package   []string `json:"package"`
	AppSecret string   `yaml:"app_secret"`
	MaxRetry  int      `yaml:"max_retry"`
}

type MeizuAppConfig struct {
	Enabled  bool   `yaml:"enabled"`
	AppID    string `yaml:"app_id"`
	AppKey   string `yaml:"app_key"`
	MaxRetry int    `yaml:"max_retry"`
}

type Config struct {
	HTTP    HTTPConfig         `yaml:"http"`
	GRPC    GRPCConfig         `yaml:"grpc"`
	IOS     []iOSAppConfig     `yaml:"ios"`
	Huawei  []HuaweiAppConfig  `yaml:"huawei"`
	Android []AndroidAppConfig `yaml:" android"`
	Vivo    []VivoAppConfig    `yaml:"vivo"`
	Oppo    []OppoAppConfig    `yaml:"oppo"`
	Xiaomi  []XiaomiAppConfig  `yaml:"xiaomi"`
	Meizu   []MeizuAppConfig   `yaml:"meizu"`
}

type HTTPConfig struct {
	Enabled bool `yaml:"enabled"`

	Address string `mapstructure:"address" yaml:"address"`
	Port    int    `mapstructure:"port" yaml:"port"`
}

func (c HTTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type GRPCConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `mapstructure:"address" yaml:"address"`
	Port    int    `mapstructure:"port" yaml:"port"`
}

func (c GRPCConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
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
