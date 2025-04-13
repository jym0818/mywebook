package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jym/mywebook/internal/service"
	"github.com/jym/mywebook/internal/service/oauth2/wechat"
	ijwt "github.com/jym/mywebook/internal/web/jwt"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	stateKey []byte
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, jwtHdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("95osj3fUD7foxmlYdDbncXz4VD2igvf1"),
		Handler:  jwtHdl,
	}
}
func (h *OAuth2WechatHandler) RegisterRouters(s *gin.Engine) {
	og := s.Group("/oauth2/wechat")
	og.GET("/authurl", h.Authurl)
	og.Any("/callback", h.Callback)

}

func (h *OAuth2WechatHandler) Authurl(ctx *gin.Context) {
	state := uuid.NewString()
	url, err := h.svc.Authurl(ctx.Request.Context(), state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "构建微信登录URL错误",
		})
		return
	}
	// 设置jwt-state
	err = h.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	err := h.verifyState(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}

	info, err := h.svc.VerifyCode(ctx.Request.Context(), code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}

	//保持登录---查找用户信息/创建用户信息
	user, err := h.userSvc.FindOrCreateByWechat(ctx.Request.Context(), info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}
	//jwt保持登录--把公共jwt提取出来
	err = h.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}

func (h *OAuth2WechatHandler) setStateCookie(c *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		StateClaims{
			State: state,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
			},
		})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	c.SetCookie("jwt-state", tokenStr,
		600, "/oauth2/wechat/callback",
		"", true, true)
	return nil
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context, state string) error {
	tokenStr, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到 state 的 cookie, %w", err)
	}
	var claims StateClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token 已经过期了, %w", err)
	}
	if claims.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
