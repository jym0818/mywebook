package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetToken(ctx *gin.Context, uid int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	ExtractToken(ctx *gin.Context) string
	CheckSession(ctx *gin.Context, ssid string) bool
}
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
	Ssid      string
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}
