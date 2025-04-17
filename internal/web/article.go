package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/service"
	ijwt "github.com/jym/mywebook/internal/web/jwt"
	"net/http"
	"strconv"
	"time"
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
	g.GET("/detail/:id", h.Detail)
	pub := g.Group("/pub")
	//pub.GET("/pub", a.PubList)
	pub.GET("/:id", h.PubDetail)
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

func (h *ArticleHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}
	claims := c.MustGet("claims").(ijwt.UserClaims)
	art, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

		return
	}
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != claims.Uid {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。

		return
	}
	c.JSON(http.StatusOK, Result{
		Data: ArticleVO{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: art.Ctime.Format(time.DateTime),
			Utime: art.Utime.Format(time.DateTime),
		},
	})

}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})

		return
	}
	art, err := h.svc.GetPublishedById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 要把作者信息带出去
			AuthorName: art.Author.Name,
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
		},
	})
}
