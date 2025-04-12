package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/service"
	"github.com/jym/mywebook/internal/service/oauth2/wechat"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}
func (h *OAuth2WechatHandler) RegisterRouters(s *gin.Engine) {
	og := s.Group("/oauth2/wechat")
	og.GET("/authurl", h.Authurl)
	og.Any("/callback", h.Callback)

}

func (h *OAuth2WechatHandler) Authurl(ctx *gin.Context) {
	url, err := h.svc.Authurl(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "构建微信登录URL错误",
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
	info, err := h.svc.VerifyCode(ctx.Request.Context(), code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}
	//保持登录---查找用户信息/创建用户信息
	user, err := h.userSvc.FindOrCreateByWechat(ctx.Request.Context(), info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}
	//jwt保持登录--把公共jwt提取出来
	err = h.setJWT(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}
