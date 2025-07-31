package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"webook/feed/events"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	sg := sarama.NewConfig()
	sg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addrs, sg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumer(feedConsumer *events.FeedEventConsumer) []events.Consumer {
	return []events.Consumer{feedConsumer}
}
