package ioc

import (
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	resolver2 "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "webook/api/proto/gen/intr/v1"
)

func InitIntrClient(client *etcdv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Target string `yaml:"target"`
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.interactive", &cfg)
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
	remote := intrv1.NewInteractiveServiceClient(cc)
	return remote
}
