package ioc

import (
	"github.com/spf13/viper"
	"webook/oauth2/service"
	"webook/pkg/logger"
)

func InitService(l logger.LoggerV1) service.Oauth2Service {
	type Config struct {
		AppID     string `yaml:"appId"`
		AppSecret string `yaml:"appSecret"`
	}
	var cfg Config
	err := viper.UnmarshalKey("weChatConf", &cfg)
	if err != nil {
		panic(err)
	}
	return service.NewService(cfg.AppID, cfg.AppSecret, l)
}
