package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "webook/mysqljob/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
)

func NewGrpcxServer(jobSvc *grpc2.MysqlCronJobServiceServer, l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
	)
	// 这里是我们封装的反向去调
	jobSvc.Register(s)

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
