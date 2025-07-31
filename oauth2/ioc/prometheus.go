package ioc

import (
	"github.com/spf13/viper"
	"webook/oauth2/service"
	"webook/pkg/logger"
)

func InitPrometheus(l logger.LoggerV1) service.Oauth2Service {
	svc := InitService(l)
	type Config struct {
		NameSpace  string `yaml:"nameSpace"`
		Subsystem  string `yaml:"subsystem"`
		InstanceID string `yaml:"instanceId"`
		Name       string `yaml:"name"`
	}
	var cfg Config
	err := viper.UnmarshalKey("prometheus", &cfg)
	if err != nil {
		panic(err)
	}
	return service.NewPrometheusDecorator(svc, cfg.NameSpace, cfg.Subsystem, cfg.InstanceID, cfg.Name)
}
