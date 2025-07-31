package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

func main() {
	InitViper()
	app := InitApp()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
	for _, consumer := range app.consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}
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
