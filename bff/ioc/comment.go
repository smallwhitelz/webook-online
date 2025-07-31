package ioc

import (
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	resolver2 "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	commentv1 "webook/api/proto/gen/comment/v1"
)

func InitCommentClient(client *etcdv3.Client) commentv1.CommentServiceClient {
	type Config struct {
		Target string `yaml:"target"`
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.comment", &cfg)
	if err != nil {
		panic(err)
	}
	resolver, err := resolver2.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(resolver)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	remote := commentv1.NewCommentServiceClient(cc)
	return remote
}
