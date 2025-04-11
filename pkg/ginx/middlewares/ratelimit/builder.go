package ratelimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/pkg/ratelimit"
	"net/http"
)

type Builder struct {
	prefix  string
	limiter ratelimit.Limiter
}

func NewBuilder(limiter ratelimit.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: limiter,
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
	return b.limiter.Limit(c.Request.Context(), key)
}
