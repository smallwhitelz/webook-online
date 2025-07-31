//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/bff/ioc"
	"webook/bff/web"
	"webook/bff/web/jwt"
)

func InitApp() *App {
	wire.Build(
		ioc.InitEtcd,
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitGinMiddlewares,
		ioc.InitUserClient,
		ioc.InitCodeClient,
		ioc.InitOauth2Client,
		ioc.InitArticleClient,
		ioc.InitIntrClient,
		ioc.InitCommentClient,
		ioc.InitRewardClient,
		ioc.InitWebServer,
		jwt.NewRedisJWTHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
