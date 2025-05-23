//go:build wireinject

package main

import (
	"github.com/google/wire"
	article2 "github.com/jym/mywebook/internal/events/article"
	"github.com/jym/mywebook/internal/repository"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/jym/mywebook/internal/repository/dao"
	"github.com/jym/mywebook/internal/service"
	"github.com/jym/mywebook/internal/web"
	ijwt "github.com/jym/mywebook/internal/web/jwt"
	"github.com/jym/mywebook/ioc"
)

var RankingService = wire.NewSet(
	repository.NewrankingRepository,
	cache.NewRedisRankingCache,
	service.NewBatchRankingService,
)

var UserService = wire.NewSet(
	dao.NewuserDAO,
	cache.NewRedisUserCache,
	repository.NewuserRepository,
	service.NewuserService,
	web.NewUserHandler,
)

var CodeService = wire.NewSet(

	cache.NewRedisCodeCache,
	repository.NewcodeRepository,
	service.NewcodeService,
)

var ArticleService = wire.NewSet(
	dao.NewarticleDAO,
	repository.NewarticleRepository,
	service.NewarticleService,
	cache.NewRedisArticle,
)

func InitWeb() *App {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitSMS,
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.InitLogger,

		ioc.InitIntrGRPCClient,
		ioc.InitEtcd,

		UserService,
		CodeService,
		ioc.InitWechat,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJwt,
		ArticleService,
		web.NewArticleHandler,

		ioc.InitKafka,
		ioc.NewSyncProducer,

		article2.NewKafkaProducer,

		ioc.InitJobs,
		ioc.InitRankingJob,
		RankingService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
