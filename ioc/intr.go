package ioc

import (
	"github.com/fsnotify/fsnotify"
	intrv1 "github.com/jym/mywebook/api/proto/gen/intr/v1"
	"github.com/jym/mywebook/interactive/service"
	"github.com/jym/mywebook/internal/web/client"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrGRPCClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
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
