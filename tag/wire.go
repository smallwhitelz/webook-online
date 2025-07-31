//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/tag/events"
	"webook/tag/grpc"
	"webook/tag/ioc"
	"webook/tag/repository"
	"webook/tag/repository/cache"
	"webook/tag/repository/dao"
	"webook/tag/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitLogger, ioc.InitKafka, ioc.InitRedis, ioc.InitSaramaSyncProducer,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		ioc.InitGrpc,
		events.NewSaramaSyncProducer,
		grpc.NewTagServiceServer,
		cache.NewTagRedisCache,
		dao.NewTagGORMDAO,
		repository.NewTagRepository,
		service.NewTagService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
