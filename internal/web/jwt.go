package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

type jwtHandler struct {
	refreshTokenKey []byte
	accessTokenKey  []byte
}

func NewJwtHandler() jwtHandler {
	return jwtHandler{
		refreshTokenKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		accessTokenKey:  []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx"),
	}
}

func (h jwtHandler) setLoginJWT(c *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setJWT(c, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshJWT(c, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (h jwtHandler) setJWT(c *gin.Context, uid int64, ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		Uid:       uid,
		UserAgent: c.GetHeader("User-Agent"),
		Ssid:      ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	tokenStr, err := token.SignedString(h.accessTokenKey)

	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (h jwtHandler) setRefreshJWT(c *gin.Context, uid int64, ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})
	tokenStr, err := token.SignedString(h.refreshTokenKey)

	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func ExtractToken(c *gin.Context) string {
	tokenStr := c.GetHeader("Authorization")
	sges := strings.SplitN(tokenStr, " ", 2)
	//传的格式错误，瞎几把传的，相当于没登陆
	return sges[1]
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
