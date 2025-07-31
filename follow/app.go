package main

import (
	"webook/follow/events"
	"webook/pkg/grpcx"
)

type App struct {
	server   *grpcx.Server
	consumer []events.Consumer
}
