//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/feed/events"
	"webook/feed/grpc"
	"webook/feed/ioc"
	"webook/feed/repository"
	"webook/feed/repository/dao"
	"webook/feed/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitLogger, ioc.InitKafka, ioc.InitEtcd,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		ioc.InitConsumer,
		events.NewFeedEventConsumer,
		ioc.InitGrpc,
		grpc.NewFeedServiceServer,
		ioc.InitFollowClient,
		ioc.RegisterHandler,
		dao.NewFeedPullEventGORMDAO,
		dao.NewFeedPushEventGORMDAO,
		repository.NewFeedEventRepo,
		service.NewFeedService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
