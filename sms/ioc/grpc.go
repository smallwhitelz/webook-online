package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
	grpc2 "webook/sms/grpc"
)

func NewGrpcxServer(smsSvc *grpc2.SmsServiceServer, l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	s := grpc.NewServer()
	// 这里是我们封装的反向去调
	smsSvc.Register(s)

	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	return &grpcx.Server{
		Server:   s,
		EtcdAddr: cfg.EtcdAddr,
		Port:     cfg.Port,
		Name:     cfg.Name,
		L:        l,
	}
}
