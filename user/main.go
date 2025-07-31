package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	InitViper()
	initPrometheus()
	app := InitApp()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}

func InitViper() {
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

func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.Handle("/user-metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
