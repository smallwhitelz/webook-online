package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "webook/code/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
)

func NewGrpcxServer(codeSvc *grpc2.CodeServiceServer, l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	s := grpc.NewServer()
	// 这里是我们封装的反向去调
	codeSvc.Register(s)

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
