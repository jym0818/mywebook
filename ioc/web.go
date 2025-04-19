package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/web"
	ijwt "github.com/jym/mywebook/internal/web/jwt"
	"github.com/jym/mywebook/internal/web/middlewares"
	"github.com/jym/mywebook/pkg/ginx/middlewares/metric"
	"github.com/jym/mywebook/pkg/ginx/middlewares/ratelimit"
	ratelimit2 "github.com/jym/mywebook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, userHandler *web.UserHandler, OAuth2WechatHandler *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	r := gin.Default()
	r.Use(mdls...)
	userHandler.RegisterRouters(r)
	OAuth2WechatHandler.RegisterRouters(r)
	articleHdl.RegisterRouters(r)
	return r
}
func InitMiddlewares(cmd redis.Cmdable, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(ratelimit2.NewRedisSlideWindowLimiter(cmd, time.Second, 100)).Build(),
		cors.New(cors.Config{
			// 允许的源地址（CORS中的Access-Control-Allow-Origin）
			// AllowOrigins: []string{"https://foo.com"},
			// 允许的 HTTP 方法（CORS中的Access-Control-Allow-Methods）
			//如果省略，那么所有方法都允许
			AllowMethods: []string{"PUT", "PATCH"},
			// 允许的 HTTP 头部（CORS中的Access-Control-Allow-Headers）
			AllowHeaders: []string{"Origin"},
			// 暴露的 HTTP 头部（CORS中的Access-Control-Expose-Headers）
			ExposeHeaders: []string{"Content-Length", "x-jwt-token", "x-refresh-token"},
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
		}),
		(&metric.MiddlewareBuilder{
			Namespace:  "jym",
			Subsystem:  "webook",
			Name:       "gin_http",
			Help:       "统计gin的http接口",
			InstanceID: "1",
		}).Build(),
		otelgin.Middleware("webook"),
		middlewares.NewLoginMiddlewareBuilder(jwtHdl).
			IgnorePath("/user/login").IgnorePath("/user/signup").IgnorePath("/user/sms/send_code").
			IgnorePath("/user/sms/login_sms").IgnorePath("/oauth2/wechat/authurl").IgnorePath("/oauth2/wechat/callback").
			IgnorePath("/user/refresh_token").
			Build(),
	}

}
