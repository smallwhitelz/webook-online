package main

import (
	"webook/feed/events"
	"webook/pkg/grpcx"
)

type App struct {
	server   *grpcx.Server
	consumer []events.Consumer
}
