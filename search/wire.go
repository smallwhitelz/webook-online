//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/search/events"
	"webook/search/grpc"
	"webook/search/ioc"
	"webook/search/repository"
	"webook/search/repository/dao"
	"webook/search/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitES, ioc.InitSamara, ioc.InitLogger)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		ioc.NewConsumers,
		ioc.InitGrpc,
		events.NewSyncDataEventConsumer,
		events.NewArticleConsumer,
		events.NewInteractiveConsumer,
		events.NewUserConsumer,
		grpc.NewSearchServiceServer,
		grpc.NewSyncServiceServer,
		dao.NewAnyESDAO,
		dao.NewArticleESDAO,
		dao.NewCollectESDAO,
		dao.NewLikeESDAO,
		dao.NewTagESDAO,
		dao.NewUserESDAO,
		repository.NewAnyRepository,
		repository.NewArticleSearchRepository,
		repository.NewUserSearchRepository,
		service.NewSyncService,
		service.NewSearchService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
