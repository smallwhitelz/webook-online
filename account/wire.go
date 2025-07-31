//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/account/grpc"
	"webook/account/ioc"
	"webook/account/repository"
	"webook/account/repository/cache"
	"webook/account/repository/dao"
	"webook/account/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitLogger, ioc.InitRedis)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		dao.NewAccountGORMDAO,
		cache.NewAccountRedisCache,
		repository.NewAccountRepository,
		service.NewAccountService,
		ioc.NewGrpcxServer,
		grpc.NewAccountServiceServer,
		wire.Struct(new(App), "*"))
	return new(App)
}
