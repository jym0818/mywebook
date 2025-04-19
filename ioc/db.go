package ioc

import (
	"github.com/jym/mywebook/internal/repository/dao"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"time"
)

func InitDB() *gorm.DB {

	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("mysql", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	//prometheus
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Namespace", "Subsystem", "Name"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	//监控sql查询时间
	vector := prometheus2.NewSummaryVec(prometheus2.SummaryOpts{
		Namespace: "jym",
		Subsystem: "webook",
		Name:      "gorm_sql_time",
		Help:      "统计sql耗时",
		Objectives: map[float64]float64{
			0.5: 0.01, 0.9: 0.01, 0.99: 0.001,
		},
	}, []string{"type", "table"})
	err = prometheus2.Register(vector)
	if err != nil {
		panic(err)
	}

	//create
	err = db.Callback().Create().Before("*").Register("prometheus_create_before", func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	})
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").Register("prometheus_create_after", func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			//啥都干不了
			return
		}
		//准备上报prometheus
		vector.WithLabelValues("create", db.Statement.Table).Observe(float64(time.Since(startTime).Milliseconds()))

	})
	if err != nil {
		panic(err)
	}
	//tracing
	db.Use(tracing.NewPlugin(tracing.WithDBName("webook"),
		//不要记录metric  我们使用了prometheus
		tracing.WithoutMetrics(),
		//不要记录查询参数，安全需求线上不要记录
		tracing.WithoutQueryVariables(),
	))

	return db
}

type Callbacks struct {
}
