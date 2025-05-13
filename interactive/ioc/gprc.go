package ioc

import (
	grpc2 "github.com/jym/mywebook/interactive/grpc"
	"github.com/jym/mywebook/pkg/grpcx"
	"github.com/jym/mywebook/pkg/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServe(intr *grpc2.InteractiveServiceServer, l logger.Logger) *grpcx.Server {
	type Config struct {
		Port      int      `yaml:"port"`
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}

	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	intr.Register(s)
	return &grpcx.Server{
		Server:    s,
		Port:      cfg.Port,
		EtcdAddrs: cfg.EtcdAddrs,
		Name:      "interactive",
		L:         l,
	}
}
