//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/code/grpc"
	"webook/code/ioc"
	"webook/code/repository"
	"webook/code/repository/cache"
	"webook/code/service"
)

func InitApp() *App {
	wire.Build(
		ioc.InitRedis,
		ioc.InitEtcd,
		ioc.InitLogger,
		ioc.InitSmslient,
		ioc.NewGrpcxServer,
		grpc.NewCodeServiceServer,
		cache.NewRedisCodeCache,
		repository.NewCodeRepository,
		service.NewCodeService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
