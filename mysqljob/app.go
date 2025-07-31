package main

import (
	"webook/mysqljob/job"
	"webook/pkg/grpcx"
)

type App struct {
	server *grpcx.Server
	job    *job.Scheduler
}
