package main

import (
	"github.com/robfig/cron/v3"
	"webook/pkg/grpcx"
)

type App struct {
	server *grpcx.Server
	cron   *cron.Cron
}
