package ioc

import (
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"webook/search/repository/dao"
)

// InitES 初始化ES客户端
func InitES() *elastic.Client {
	type Config struct {
		Url string `yaml:"url"`
	}
	var cfg Config
	err := viper.UnmarshalKey("es", &cfg)
	if err != nil {
		panic(err)
	}
	opts := []elastic.ClientOptionFunc{
		elastic.SetSniff(false),
		elastic.SetURL(cfg.Url),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	err = dao.InitES(client)
	if err != nil {
		panic(err)
	}
	return client
}
