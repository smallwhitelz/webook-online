package ioc

import (
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	resolver2 "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	articlev1 "webook/api/proto/gen/article/v1"
)

func InitArticleClient(client *etcdv3.Client) articlev1.ArticleServiceClient {
	type Config struct {
		Target string `yaml:"target"`
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.article", &cfg)
	if err != nil {
		panic(err)
	}
	resolver, err := resolver2.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{
		grpc.WithResolvers(resolver),
	}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	remote := articlev1.NewArticleServiceClient(cc)
	return remote
}
