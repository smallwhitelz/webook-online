package main

import (
	"webook/pkg/grpcx"
	"webook/reward/events"
)

type App struct {
	server   *grpcx.Server
	consumer events.Consumer
}
