package ioc

import (
	grpc2 "github.com/jym/mywebook/interactive/grpc"
	"github.com/jym/mywebook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServe(intr *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}

	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	intr.Register(s)
	return &grpcx.Server{
		Addr:   cfg.Addr,
		Server: s,
	}
}
