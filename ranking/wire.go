//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/ranking/grpc"
	"webook/ranking/ioc"
	"webook/ranking/repository"
	"webook/ranking/repository/cache"
	"webook/ranking/service"
)

func InitApp() *App {
	wire.Build(
		ioc.InitRedis,
		ioc.InitRlockClient,
		ioc.InitEtcd,
		ioc.InitLogger,
		ioc.InitArtClient,
		ioc.InitIntrClient,
		ioc.InitJobs,
		ioc.InitRankingJob,
		ioc.NewGrpcxServer,
		grpc.NewRankingServiceServer,
		cache.NewRankingRedisCache,
		repository.NewCachedRankingRepository,
		service.NewBatchRankingService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
