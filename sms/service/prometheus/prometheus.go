package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"webook/sms/service"
)

// Decorator 装饰器模式实现对sms进行监控
type Decorator struct {
	svc    service.SmsService
	vector *prometheus.SummaryVec
}

func NewDecorator(svc service.SmsService, opt prometheus.SummaryOpts) *Decorator {
	return &Decorator{
		svc:    svc,
		vector: prometheus.NewSummaryVec(opt, []string{"tpl_id"}),
	}

}

func (d *Decorator) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	start := time.Now()
	// 监控响应时间
	defer func() {
		duration := time.Since(start).Milliseconds()
		d.vector.WithLabelValues(tplId).Observe(float64(duration))
	}()
	return d.svc.Send(ctx, tplId, args, numbers...)
}
