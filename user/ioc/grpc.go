package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"webook/pkg/grpcx"
	logger2 "webook/pkg/grpcx/interceptor/logger"
	"webook/pkg/grpcx/interceptor/prometheus"
	"webook/pkg/grpcx/interceptor/trace"
	grpc3 "webook/pkg/grpcx/monitor"
	"webook/pkg/logger"
	grpc2 "webook/user/grpc"
)

func NewGrpcxServer(userSvc *grpc2.UserServiceServer, l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	grpc3.InitZipkin("user")
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// 接入服务端日志
			logger2.NewInterceptorBuilder(l).BuildServerUnaryInterceptor(),
			// 接入服务端trace
			trace.NewOTELInterceptorBuilder("server_user", nil, nil).
				BuildUnaryServerInterceptor(),
			prometheus.NewInterceptorBuilder("grpc_webook", "user", "req").BuildServerUnaryInterceptor(),
		),
	)
	// 这里是我们封装的反向去调
	userSvc.Register(s)

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
