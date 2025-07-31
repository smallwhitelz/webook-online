//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/mysqljob/grpc"
	"webook/mysqljob/ioc"
	"webook/mysqljob/repository"
	"webook/mysqljob/repository/dao"
	"webook/mysqljob/service"
)

func InitApp() *App {
	wire.Build(
		ioc.InitDB,
		ioc.InitLogger,
		ioc.NewGrpcxServer,
		grpc.NewMysqlCronJobServiceServer,
		dao.NewGORMJobDAO,
		repository.NewPreemptJobRepository,
		service.NewCronJobService,
		ioc.InitScheduler,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
