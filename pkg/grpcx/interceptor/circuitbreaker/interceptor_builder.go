package circuitbreaker

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InterceptorBuilder 熔断
// 这里熔断我们就借助Kratos里的熔断功能
type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

func (b *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		err = b.breaker.Allow()
		if err == nil {
			resp, err = handler(ctx, req)
			if err == nil {
				b.breaker.MarkSuccess()
			} else {
				// 可以更加仔细的检测，只有真实代表服务端出现故障的，才 mark failed
				// MarkFailed的前提一定是系统出现故障而不是业务出现故障
				// 什么叫业务故障？例如你要求参数传入一个大于0的，结果用户传入小于0，这就是业务故障
				b.breaker.MarkFailed()
			}
			return
		} else {
			// 触发了熔断
			b.breaker.MarkFailed()
			return nil, status.Errorf(codes.Unavailable, "熔断")
		}
	}
}
