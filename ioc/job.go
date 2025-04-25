package ioc

import (
	"fmt"
	"github.com/jym/mywebook/internal/job"
	"github.com/jym/mywebook/internal/service"
	"github.com/jym/mywebook/pkg/logger"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService) job.Job {
	return job.NewRankingJob(svc)
}

func InitJobs(l logger.Logger, rankingJob job.Job) *cron.Cron {
	c := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	_, err := c.AddJob("0 */3 * * * ?", cbd.Build(rankingJob))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return c
}
