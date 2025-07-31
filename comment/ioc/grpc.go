package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "webook/comment/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
)

func NewGrpcServer(cmtSvc *grpc2.CommentServiceServer, l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	server := grpc.NewServer()
	cmtSvc.Register(server)
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	return &grpcx.Server{
		Server:   server,
		EtcdAddr: cfg.EtcdAddr,
		Port:     cfg.Port,
		Name:     cfg.Name,
		L:        l,
	}
}
