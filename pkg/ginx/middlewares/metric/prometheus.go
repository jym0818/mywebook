package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceID string
}

func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_resp_time",
		Help:      m.Help,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
		Objectives: map[float64]float64{0.5: 0.01, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
	}, []string{"method", "pattern", "status"})
	prometheus.MustRegister(summary)
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			duartion := time.Since(start)

			pattern := c.FullPath()
			//如果是404呢 方便表示
			if pattern == "" {
				pattern = "unknown"
			}
			summary.WithLabelValues(c.Request.Method, pattern, strconv.Itoa(c.Writer.Status())).Observe(float64(duartion.Milliseconds()))
		}()
		c.Next()

	}
}
