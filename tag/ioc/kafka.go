package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
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
	// 这里要注意顺序
	// 一个用户他可能给一个资源打标签，先打abc，再打的bcd
	// 这里要注意消息发送的顺序，所以要用hash类partition
	sg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	client, err := sarama.NewClient(cfg.Addrs, sg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}
