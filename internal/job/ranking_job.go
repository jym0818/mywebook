package job

import (
	"context"
	"github.com/jym/mywebook/internal/service"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func NewRankingJob(svc service.RankingService) *RankingJob {
	//保证过期时间内计算完数据
	return &RankingJob{svc: svc, timeout: time.Minute}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}
