package job

import (
	"context"
	"github.com/jym/mywebook/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	p      *prometheus.SummaryVec
	l      logger.Logger
	tracer trace.Tracer
}

func NewCronJobBuilder(l logger.Logger) *CronJobBuilder {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "jym",
		Subsystem: "webook",
		Name:      "job",
		Help:      "cron运行",
	}, []string{"job_name", "success"})
	prometheus.MustRegister(summary)
	return &CronJobBuilder{
		l:      l,
		p:      summary,
		tracer: otel.GetTracerProvider().Tracer("github.com/jym/webook/internal/job/job_builder.go"),
	}
}

func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobFunAdapter(func() error {
		_, span := b.tracer.Start(context.Background(), name)
		defer span.End()
		start := time.Now()
		b.l.Info("start job run", logger.String("job_name", name))
		var success bool
		defer func() {
			duartion := time.Since(start).Milliseconds()
			b.p.WithLabelValues(name, strconv.FormatBool(success)).Observe(float64(duartion))
		}()
		err := job.Run()
		success = err == nil
		if err != nil {
			span.RecordError(err)
			b.l.Info("stop job run", logger.String("job_name", name))
			b.l.Error(err.Error(), logger.String("job_name", name))
		}
		return nil
	})
}

type cronJobFunAdapter func() error

func (c cronJobFunAdapter) Run() {
	_ = c()
}
