package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/repository"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/jym/mywebook/internal/repository/dao"
	"github.com/jym/mywebook/internal/service"
	"github.com/jym/mywebook/internal/service/sms/memory"
	"github.com/jym/mywebook/internal/web"
	"github.com/jym/mywebook/internal/web/middlewares"
	"github.com/jym/mywebook/pkg/ginx/middlewares/ratelimit"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	r := gin.Default()
	cmd := InitRedis()
	// 一分钟 100 次。
	r.Use(ratelimit.NewBuilder(cmd, time.Minute, 100).Build())

	r.Use(cors.New(cors.Config{
		// 允许的源地址（CORS中的Access-Control-Allow-Origin）
		// AllowOrigins: []string{"https://foo.com"},
		// 允许的 HTTP 方法（CORS中的Access-Control-Allow-Methods）
		//如果省略，那么所有方法都允许
		AllowMethods: []string{"PUT", "PATCH"},
		// 允许的 HTTP 头部（CORS中的Access-Control-Allow-Headers）
		AllowHeaders: []string{"Origin"},
		// 暴露的 HTTP 头部（CORS中的Access-Control-Expose-Headers）
		ExposeHeaders: []string{"Content-Length", "x-jwt-token"},
		// 是否允许携带身份凭证（CORS中的Access-Control-Allow-Credentials）
		AllowCredentials: true,
		// 允许源的自定义判断函数，返回true表示允许，false表示不允许
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 允许你的开发环境
				return true
			}
			// 允许包含 "yourcompany.com" 的源
			return strings.Contains(origin, "yourcompany.com")
		},
		// 用于缓存预检请求结果的最大时间（CORS中的Access-Control-Max-Age）
		MaxAge: 12 * time.Hour,
	}))
	r.Use(middlewares.NewLoginMiddlewareBuilder().
		IgnorePath("/user/login").IgnorePath("/user/signup").IgnorePath("/user/sms/send_code").
		IgnorePath("/user/sms/login_sms").
		Build())
	//依赖注入  解耦
	db, err := gorm.Open(mysql.Open("root:root@tcp(47.120.51.5:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	userDAO := dao.NewuserDAO(db)
	userCache := cache.NewRedisUserCache(cmd)
	repo := repository.NewuserRepository(userDAO, userCache)
	svc := service.NewuserService(repo)
	//codeSvc
	smsSvc := memory.NewService()
	codeCache := cache.NewRedisCodeCache(cmd)
	codeRepo := repository.NewcodeRepository(codeCache)
	codeSvc := service.NewcodeService(codeRepo, smsSvc)

	user := web.NewUserHandler(svc, codeSvc)
	user.RegisterRouters(r)
	r.Run(":8080")
}

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "47.120.51.5:6379",
		Password: "",
		DB:       1,
	})
	return cmd
}
