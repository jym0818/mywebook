package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	RefreshTokenKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	AccessTokenKey  = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

type RedisJwt struct {
	cmd redis.Cmdable
}

func NewRedisJwt(cmd redis.Cmdable) Handler {
	return &RedisJwt{
		cmd: cmd,
	}
}

func (r *RedisJwt) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.setRefreshJWT(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisJwt) SetToken(ctx *gin.Context, uid int64, ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		Ssid:      ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	tokenStr, err := token.SignedString(AccessTokenKey)

	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r *RedisJwt) ClearToken(ctx *gin.Context) error {
	claims := ctx.MustGet("claims").(UserClaims)
	err := r.cmd.Set(ctx, fmt.Sprintf("user:ssid:%s", claims.Ssid), "true", time.Hour*24*7).Err()
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-refresh", "")
	return nil
}

func (r *RedisJwt) ExtractToken(ctx *gin.Context) string {
	tokenStr := ctx.GetHeader("Authorization")
	sges := strings.SplitN(tokenStr, " ", 2)
	//传的格式错误，瞎几把传的，相当于没登陆
	return sges[1]
}

func (r *RedisJwt) CheckSession(ctx *gin.Context, ssid string) bool {
	cnt, err := r.cmd.Exists(ctx, fmt.Sprintf("user:ssid:%s", ssid)).Result()
	if err != nil || cnt > 0 {
		return false
	}
	return true
}
func (r *RedisJwt) setRefreshJWT(c *gin.Context, uid int64, ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})
	tokenStr, err := token.SignedString(RefreshTokenKey)

	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}
