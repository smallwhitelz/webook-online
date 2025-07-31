//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"webook/mysqljob/job"
	"webook/mysqljob/repository"
	"webook/mysqljob/repository/dao"
	"webook/mysqljob/service"
)

var jobProviderSet = wire.NewSet(
	service.NewCronJobService,
	repository.NewPreemptJobRepository,
	dao.NewGORMJobDAO,
)

func InitJobScheduler() *job.Scheduler {
	wire.Build(
		InitDB,
		InitLogger,
		jobProviderSet,
		job.NewScheduler)
	return &job.Scheduler{}
}
