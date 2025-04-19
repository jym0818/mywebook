package metric

import (
	"context"
	"github.com/jym/mywebook/internal/service/sms"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type PrometheusService struct {
	svc    sms.Service
	vector *prometheus.SummaryVec
}

func NewPrometheusService(svc sms.Service) sms.Service {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "jym",
		Subsystem:  "webook",
		Name:       "sms",
		Help:       "发送短信耗时",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"tpl"})
	prometheus.MustRegister(vector)
	return &PrometheusService{
		svc:    svc,
		vector: vector,
	}
}
func (p *PrometheusService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	start := time.Now()
	defer func() {
		p.vector.WithLabelValues(tpl).Observe(float64(time.Since(start).Milliseconds()))
	}()

	return p.svc.Send(ctx, tpl, args, numbers...)

}
