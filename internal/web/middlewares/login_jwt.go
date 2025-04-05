package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym/mywebook/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
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
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		sges := strings.SplitN(tokenStr, " ", 2)
		//传的格式错误，瞎几把传的，相当于没登陆
		if len(sges) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var claims web.UserClaims
		token, err := jwt.ParseWithClaims(sges[1], &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid || token == nil || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//刷新登录
		if claims.ExpiresAt.Time.Sub(time.Now()) < 45*time.Second {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"))
			if err != nil {
			}
			c.Header("x-jwt-token", tokenStr)
		}

		c.Set("claims", claims)
	}
}
