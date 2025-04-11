package ioc

import (
	"github.com/jym/mywebook/internal/service/sms"
	"github.com/jym/mywebook/internal/service/sms/failover"
	"github.com/jym/mywebook/internal/service/sms/memory"
	"github.com/jym/mywebook/internal/service/sms/ratelimit"
	ratelimit2 "github.com/jym/mywebook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSMS(cmd redis.Cmdable) sms.Service {
	svc := memory.NewService()
	failoverSvc := failover.NewTimeoutFailoverSMSService([]sms.Service{svc})
	limiter := ratelimit2.NewRedisSlideWindowLimiter(cmd, time.Second, 100)
	return ratelimit.NewRateLimitSMSService(failoverSvc, limiter)
}
