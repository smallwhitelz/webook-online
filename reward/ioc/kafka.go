package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"webook/reward/events"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	sg := sarama.NewConfig()
	sg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, sg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumers(pc *events.PaymentEventConsumer) events.Consumer {
	return pc
}
