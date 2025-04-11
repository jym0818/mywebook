package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/service/oauth2/wechat"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
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

func (h *OAuth2WechatHandler) Callback(context *gin.Context) {

}
