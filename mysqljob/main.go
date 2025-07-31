package main

import (
	"context"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	initViper()
	app := InitApp()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := app.job.Schedule(ctx)
		if err != nil && err != context.Canceled {
			log.Println("调度器异常退出", err)
		} else {
			log.Println("调度器已正常停止")
		}
	}()
}

func initViper() {
	cfile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	// 这一步后 cfile才有值
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println(viper.GetString("test.key"))
}
