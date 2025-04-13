package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/service"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

const emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

const passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
const biz = "login"

type UserHandler struct {
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
	svc           service.UserService
	codeSvc       service.CodeService
	jwtHandler
	cmd redis.Cmdable
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, cmd redis.Cmdable) *UserHandler {
	return &UserHandler{
		emailRegex:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
		codeSvc:       codeSvc,
		jwtHandler:    NewJwtHandler(),
		cmd:           cmd,
	}
}

func (h *UserHandler) RegisterRouters(s *gin.Engine) {
	ug := s.Group("/user")
	ug.POST("/signup", h.Signup)
	ug.POST("/login", h.Login)
	ug.POST("/profile", h.Profile)
	ug.POST("/sms/send_code", h.SendCode)
	ug.POST("/sms/login_sms", h.LoginSMS)
	ug.POST("/refresh_token", h.RefreshToken)
	ug.POST("/logout", h.Logout)
}

func (h *UserHandler) Login(c *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {

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
	if err := h.setLoginJWT(c, u.Id); err != nil {
		c.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}

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

func (h *UserHandler) SendCode(ctx *gin.Context) {

	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {

		return
	}
	//参数校验
	err := h.codeSvc.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "发送成功",
	})
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {

		return
	}
	//参数校验一下

	//校验验证码是否正确
	ok, err := h.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误2"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}

	//验证码正确  登录或者注册并登录
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误3"})
		return
	}
	//jwt
	if err := h.setLoginJWT(ctx, u.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误4"})
		return
	}

	ctx.JSON(http.StatusOK, Result{Data: u})

}
func (h *UserHandler) RefreshToken(ctx *gin.Context) {

	tokenStr := ExtractToken(ctx)
	var claims RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return h.refreshTokenKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	cnt, err := h.cmd.Exists(ctx, fmt.Sprintf("user:ssid:%s", claims.Ssid)).Result()
	if err != nil || cnt > 0 {
		//要么有问题，要么退出登录
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.setJWT(ctx, claims.Uid, claims.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})

}

func (h *UserHandler) Logout(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(UserClaims)
	err := h.cmd.Set(ctx, fmt.Sprintf("user:ssid:%s", claims.Ssid), "true", time.Hour*24*7).Err()
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "退出登录失败",
		})
	}
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-refresh", "")
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录成功",
	})
}
