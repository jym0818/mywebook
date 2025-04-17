package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/service"
	ijwt "github.com/jym/mywebook/internal/web/jwt"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func (h *ArticleHandler) RegisterRouters(s *gin.Engine) {
	g := s.Group("/articles")
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)
	g.POST("/list", h.List)
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (h *ArticleHandler) Edit(c *gin.Context) {

	var req ArticleReq
	if err := c.Bind(&req); err != nil {
		return
	}
	//claims
	claims := c.MustGet("claims").(ijwt.UserClaims)
	// 调用 svc 的代码
	id, err := h.svc.Save(c.Request.Context(), req.toDomain(claims.Uid))
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})

}

func (h *ArticleHandler) Publish(c *gin.Context) {
	var req ArticleReq
	if err := c.Bind(&req); err != nil {
		return
	}
	//claims
	claims := c.MustGet("claims").(ijwt.UserClaims)
	id, err := h.svc.Publish(c.Request.Context(), req.toDomain(claims.Uid))
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(c *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	claims := c.MustGet("claims").(ijwt.UserClaims)
	err := h.svc.Withdraw(c.Request.Context(), domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) List(c *gin.Context) {
	var req ListReq
	if err := c.Bind(&req); err != nil {
		return
	}
	claims := c.MustGet("claims").(ijwt.UserClaims)

	res, err := h.svc.List(c.Request.Context(), claims.Uid, req.Limit, req.Offset)

	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	var arts []ArticleVO
	for _, item := range res {
		arts = append(arts, ArticleVO{
			Id:         item.Id,
			Title:      item.Title,
			Status:     item.Status.ToUint8(),
			AuthorId:   item.Author.Id,
			AuthorName: item.Author.Name,
			Ctime:      item.Ctime.Format("2006-01-02 15:04:05"),
			Utime:      item.Utime.Format("2006-01-02 15:04:05"),
			Abstract:   item.Abstract(),
		})
	}
	c.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: arts,
	})
}
