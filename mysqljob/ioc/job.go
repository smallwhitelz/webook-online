package ioc

import (
	"webook/mysqljob/job"
	"webook/mysqljob/service"
	"webook/pkg/logger"
)

// InitScheduler 如何注册一个job，这里的注册肯定要和用户前端
func InitScheduler(svc service.CronJobService, l logger.LoggerV1) *job.Scheduler {
	res := job.NewScheduler(svc, l)
	res.RegisterExecutor(initLocalFuncExecutor())
	return res
}

func initLocalFuncExecutor() *job.LocalFuncExecutor {
	localExecutor := job.NewLocalFuncExecutor()
	return localExecutor
}
