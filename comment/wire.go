//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/comment/grpc"
	"webook/comment/ioc"
	"webook/comment/repository"
	"webook/comment/repository/dao"
	"webook/comment/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitLogger)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		dao.NewCommentGORMDAO,
		repository.NewCommentRepository,
		service.NewCommentService,
		ioc.NewGrpcServer,
		grpc.NewCommentServiceServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
