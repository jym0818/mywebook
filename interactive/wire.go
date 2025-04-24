//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jym/mywebook/interactive/events"
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

func InitGRPCServer() *App {
	wire.Build(
		interactiveSvc,
		ioc.InitRedis,
		ioc.InitDB,
		ioc.InitKafka,
		ioc.InitLogger,
		events.NewKafkaConsumer,
		ioc.NewConsumers,
		grpc.NewInteractiveServiceServer,
		ioc.InitGRPCxServe,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
