//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jym/mywebook/internal/repository"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/jym/mywebook/internal/repository/dao"
	"github.com/jym/mywebook/internal/service"
	"github.com/jym/mywebook/internal/web"
	ijwt "github.com/jym/mywebook/internal/web/jwt"
	"github.com/jym/mywebook/ioc"
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
)

func InitWeb() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitSMS,
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.InitLogger,
		UserService,
		CodeService,
		ioc.InitWechat,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJwt,
		ArticleService,
		web.NewArticleHandler,
	)
	return new(gin.Engine)
}
