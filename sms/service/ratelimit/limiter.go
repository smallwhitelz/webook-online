package ratelimit

import (
	"context"
	"errors"
	"webook/pkg/limiter"
	"webook/sms/service"
)

var errLimited = errors.New("触发限流")

// RateLimitSMSService 对短信服务进行限流
type RateLimitSMSService struct {
	// 被装饰的
	svc     service.SmsService
	limiter limiter.Limiter
	key     string
}

type RateLimitSMSServiceV1 struct {
	// 被装饰的
	service.SmsService
	limiter limiter.Limiter
	key     string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewRateLimitSMSService(svc service.SmsService, l limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
		key:     "sms-limiter",
	}
}
