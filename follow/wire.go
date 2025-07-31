//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/follow/events"
	"webook/follow/grpc"
	"webook/follow/ioc"
	"webook/follow/repository"
	"webook/follow/repository/cache"
	"webook/follow/repository/dao"
	"webook/follow/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitLogger, ioc.InitRedis, ioc.InitKafka)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		ioc.NewGrpcServer,
		ioc.InitConsumer,
		events.NewMysqlBinlogConsumer,
		grpc.NewFollowRelationServiceServer,
		dao.NewFollowGROMDAO,
		cache.NewFollowRedisCache,
		repository.NewFollowRelationRepository,
		service.NewFollowRelationService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
