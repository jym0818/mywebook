package main

import (
	intrv1 "github.com/jym/mywebook/api/proto/gen/intr/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	initViper()
	server := grpc.NewServer()
	defer server.GracefulStop()
	intrSvc := InitGRPCServer()
	intrv1.RegisterInteractiveServiceServer(server, intrSvc)

	l, err := net.Listen("tcp", ":8090")

	if err != nil {
		panic(err)
	}
	//阻塞在这里
	err = server.Serve(l)
	log.Println(err)
}

func initViper() {
	viper.SetConfigName("dev")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
