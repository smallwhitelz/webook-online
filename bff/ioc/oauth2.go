package ioc

import (
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	resolver2 "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	oauth2v1 "webook/api/proto/gen/oauth2/v1"
	grpc3 "webook/pkg/grpcx/monitor"
)

func InitOauth2Client(client *etcdv3.Client) oauth2v1.Oauth2ServiceClient {
	type Config struct {
		Target string `yaml:"target"`
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.oauth2", &cfg)
	if err != nil {
		panic(err)
	}
	resolver, err := resolver2.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	grpc3.InitZipkin("user")
	opts := []grpc.DialOption{
		grpc.WithResolvers(resolver)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	remote := oauth2v1.NewOauth2ServiceClient(cc)
	return remote
}
