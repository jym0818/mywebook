//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jym/mywebook/interactive/grpc"
	"github.com/jym/mywebook/interactive/ioc"
	"github.com/jym/mywebook/interactive/repository"
	"github.com/jym/mywebook/interactive/repository/cache"
	"github.com/jym/mywebook/interactive/repository/dao"
	"github.com/jym/mywebook/interactive/service"
)

var interactiveSvc = wire.NewSet(
	service.NewinteractiveService,
	repository.NewinteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(
		interactiveSvc,
		ioc.InitRedis,
		ioc.InitDB,
		grpc.NewInteractiveServiceServer,
	)
	return new(grpc.InteractiveServiceServer)
}
