//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/user/grpc"
	"webook/user/ioc"
	"webook/user/repository"
	"webook/user/repository/cache"
	"webook/user/repository/dao"
	"webook/user/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitLogger, ioc.InitRedis,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		ioc.NewGrpcxServer,
		grpc.NewUserServiceServer,
		dao.NewUserDao,
		cache.NewUserCache,
		repository.NewCachedUserRepository,
		service.NewUserService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
