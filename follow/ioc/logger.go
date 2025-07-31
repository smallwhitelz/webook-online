package ioc

import (
	"go.uber.org/zap"
	"webook/pkg/logger"
)

func InitLogger() logger.LoggerV1 {
	zap, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	l := logger.NewZapLogger(zap)
	return l
}
