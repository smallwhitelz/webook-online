package main

import (
	"webook/pkg/grpcx"
	"webook/search/events"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}
