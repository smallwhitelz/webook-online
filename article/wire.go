//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/article/events/article"
	"webook/article/grpc"
	"webook/article/ioc"
	"webook/article/repository"
	"webook/article/repository/cache"
	"webook/article/repository/dao"
	"webook/article/service"
)

func InitApp() *App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitDB,
		ioc.InitEtcd,
		ioc.InitRedis,
		ioc.InitSaramaClient,
		ioc.InitConsumers,
		ioc.InitSyncProducer,
		ioc.NewGrpcxServer,
		ioc.InitUserClient,
		article.NewMysqlBinlogConsumer,
		article.NewSaramaSyncProducer,
		grpc.NewArticleServiceServer,
		cache.NewArticleRedisCache,
		dao.NewArticleGORMDAO,
		repository.NewCachedArticleRepository,
		service.NewArticleService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
