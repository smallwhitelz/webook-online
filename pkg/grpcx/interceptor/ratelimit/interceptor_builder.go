package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"webook/pkg/limiter"
)

type InterceptorBuilder struct {
	limiter limiter.Limiter
	key     string
}

// NewInterceptorBuilder key 1. limiter:interactive-service => 整个点赞的应用限流
func NewInterceptorBuilder(limiter limiter.Limiter, key string) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key}
}

func (b *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			// 保守做法
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
			// 激进做法就是err了，继续不拒绝请求，这里的err一般是redis可能崩溃了
		}
		if limited {
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}
		return handler(ctx, req)
	}
}

// BuildServerUnaryInterceptorV1 演示降级
// 降级一般都会和熔断或者限流结合在一起
// 一个服务能不能降级、要怎么降级，降级到什么程度都是和业务强相关的
func (b *InterceptorBuilder) BuildServerUnaryInterceptorV1() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil || limited {
			// 打一个标记位
			ctx = context.WithValue(ctx, "downgrade", "true")
		}
		return handler(ctx, req)
	}
}

// BuildServerUnaryInterceptorService 服务级别的限流
func (b *InterceptorBuilder) BuildServerUnaryInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if strings.HasPrefix(info.FullMethod, "/UserService") {
			// 这个key，limiter:UserService
			limited, err := b.limiter.Limit(ctx, b.key)
			if err != nil {
				// 保守做法
				return nil, status.Errorf(codes.ResourceExhausted, "限流")
				// 激进做法就是err了，继续不拒绝请求，这里的err一般是redis可能崩溃了
			}
			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "限流")
			}
		}
		return handler(ctx, req)
	}
}
