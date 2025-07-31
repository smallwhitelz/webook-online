package ioc

import (
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"webook/pkg/gormx"
	"webook/reward/repository/dao"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(43.154.97.245:13316)/webook_reward",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	err = db.Use(prometheus.New(prometheus.Config{
		DBName: "webook_reward",
		// 每15秒采集一次数据
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				// 参数中如果有thread_running，就会设置为label
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "geektime_zl",
		Subsystem: "webook_reward",
		Name:      "gorm_db",
		Help:      "这是一个统计 GORM 的数据库查询",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	err = db.Use(cb)
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithDBName("webook_reward")))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
