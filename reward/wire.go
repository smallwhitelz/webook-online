//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/reward/events"
	"webook/reward/grpc"
	"webook/reward/ioc"
	"webook/reward/repository"
	"webook/reward/repository/cache"
	"webook/reward/repository/dao"
	"webook/reward/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitRedis, ioc.InitLogger, ioc.InitSaramaClient, ioc.InitEtcd)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		ioc.InitAccountClient,

		ioc.InitPaymentClient,
		ioc.InitConsumers,
		events.NewPaymentEventConsumer,
		ioc.NewGrpcxServer,
		grpc.NewRewardServiceServer,
		dao.NewRewardGORMDAO,
		cache.NewRewardRedisCache,
		repository.NewRewardRepository,
		service.NewWechatNativeRewardService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
