package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym/mywebook/internal/web"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
	cmd   redis.Cmdable
}

func NewLoginMiddlewareBuilder(cmd redis.Cmdable) *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		cmd: cmd,
	}
}
func (m *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	m.paths = append(m.paths, path)
	return m
}
func (m *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range m.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		tokenStr := web.ExtractToken(c)
		var claims web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid || token == nil || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//判断agent
		if c.GetHeader("User-Agent") != claims.UserAgent {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//判断ssid是否退出登录
		cnt, err := m.cmd.Exists(c, fmt.Sprintf("user:ssid:%s", claims.Ssid)).Result()
		if err != nil || cnt > 0 {
			//要么有问题，要么退出登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("claims", claims)
	}
}
