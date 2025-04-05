package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

//go:embed slide_window.lua
var luaScript string

type Builder struct {
	prefix   string
	cmd      redis.Cmdable
	interval time.Duration
	rate     int
}

func NewBuilder(cmd redis.Cmdable, interval time.Duration, rate int) *Builder {
	return &Builder{
		cmd:      cmd,
		prefix:   "ip-limiter",
		interval: interval,
		rate:     rate,
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		limited, err := b.limit(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
func (b *Builder) limit(c *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, c.ClientIP())
	return b.cmd.Eval(c, luaScript, []string{key}, b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
