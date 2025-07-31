package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"webook/follow/events"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	sgCfg := sarama.NewConfig()
	sgCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, sgCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumer(followConsumer *events.MysqlBinlogConsumer) []events.Consumer {
	return []events.Consumer{
		followConsumer,
	}
}
