package ioc

import (
	"github.com/fsnotify/fsnotify"
	intrv1 "github.com/jym/mywebook/api/proto/gen/intr/v1"
	"github.com/jym/mywebook/interactive/service"
	"github.com/jym/mywebook/internal/web/client"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEtcd() *clientv3.Client {
	var cfg clientv3.Config
	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return cli
}

func InitIntrGRPCClient(client *clientv3.Client) intrv1.InteractiveServiceClient {

	bd, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	cc, err := grpc.Dial("etcd:///service/interactive", grpc.WithResolvers(bd), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	return remote
}

func InitIntrGRPCClientV1(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string
		Threshold int32
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	cc, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	local := client.NewInteractiveServiceAdapter(svc)
	remote := intrv1.NewInteractiveServiceClient(cc)
	res := client.NewGreyScaleInteractiveServiceClient(remote, local)
	// 我的习惯是在这里监听
	viper.OnConfigChange(func(in fsnotify.Event) {
		var cfg Config
		err = viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {
			// 你可以输出日志
		}
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
