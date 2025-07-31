package ioc

import (
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	resolver2 "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	accountv1 "webook/api/proto/gen/account/v1"
)

func InitAccountClient(client *etcdv3.Client) accountv1.AccountServiceClient {
	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.account", &cfg)
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
	cc, err := grpc.NewClient(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := accountv1.NewAccountServiceClient(cc)
	return remote
}
