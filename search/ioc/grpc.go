package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
	grpc2 "webook/search/grpc"
)

func InitGrpc(searchSvc *grpc2.SearchServiceServer, syncSvc *grpc2.SyncServiceServer, l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	s := grpc.NewServer()
	searchSvc.Register(s)
	syncSvc.Register(s)
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
