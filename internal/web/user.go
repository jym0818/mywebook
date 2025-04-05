package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/service"
	"net/http"
	"time"
)

const emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

const passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

type UserHandler struct {
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
	svc           service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		emailRegex:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
	}
}

func (h *UserHandler) RegisterRouters(s *gin.Engine) {
	ug := s.Group("/user")
	ug.POST("/signup", h.Signup)
	ug.POST("/login", h.Login)
	ug.POST("/profile", h.Profile)
}

func (h *UserHandler) Login(c *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	u, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.JSON(http.StatusOK, Result{Msg: "账号或者密码错误"})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	//保持登录状态
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		Uid: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	tokenStr, err := token.SignedString([]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"))
	if err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	c.Header("x-jwt-token", tokenStr)

	c.JSON(http.StatusOK, Result{
		Msg:  "登录成功",
		Data: u,
	})
}

func (h *UserHandler) Signup(c *gin.Context) {
	//接收参数
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	//参数校验

	ok, err := h.emailRegex.MatchString(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{Msg: "邮箱格式错误"})
		return
	}

	ok, err = h.passwordRegex.MatchString(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{Msg: "密码格式错误"})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusOK, Result{Msg: "两次密码不一致"})
		return
	}

	err = h.svc.Signup(c.Request.Context(), domain.User{
		Password: req.Password,
		Email:    req.Email,
	})
	if err == service.ErrUserDuplicateEmail {
		c.JSON(http.StatusOK, Result{Msg: "邮箱已注册"})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
	}

	c.JSON(http.StatusOK, Result{Msg: "注册成功"})

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(UserClaims)
	ctx.JSON(http.StatusOK, Result{Data: claims})

}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
