package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ijwt "github.com/jym/mywebook/internal/web/jwt"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string

	ijwt.Handler
}

func NewLoginMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		Handler: jwtHdl,
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
		tokenStr := m.ExtractToken(c)
		var claims ijwt.UserClaims
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
		ok := m.CheckSession(c, claims.Ssid)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("claims", claims)
	}
}
