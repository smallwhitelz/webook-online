//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/oauth2/grpc"
	"webook/oauth2/ioc"
)

func InitApp() *App {
	wire.Build(
		ioc.InitLogger,
		ioc.NewGrpcxServer,
		ioc.InitPrometheus,
		grpc.NewOauth2ServiceServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
