//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/sms/grpc"
	"webook/sms/ioc"
)

func InitApp() *App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitSmsService,
		ioc.NewGrpcxServer,
		grpc.NewSmsServiceServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
