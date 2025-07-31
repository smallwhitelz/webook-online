package main

import (
	"webook/article/events"
	"webook/pkg/grpcx"
)

type App struct {
	server    *grpcx.Server
	consumers []events.Consumer
}
