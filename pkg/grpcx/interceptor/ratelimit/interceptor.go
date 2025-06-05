package ratelimit

import (
	"context"
	"github.com/jym/mywebook/pkg/logger"
	"github.com/jym/mywebook/pkg/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type InterceptorBuilder struct {
	limiter ratelimit.Limiter
	l       logger.Logger
	key     string
}

func NewInterceptorBuilder(limiter ratelimit.Limiter, key string, l logger.Logger) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key, l: l}
}
func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			// err 不为nil，你要考虑你用保守的，还是用激进的策略
			// 这是保守的策略
			b.l.Error("判定限流出现问题")
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")

			// 这是激进的策略
			// return handler(ctx, req)
		}
		if limited {
			//defVal, ok := b.defaultValueMap[info.FullMethod]
			//if ok {
			//	err = json.Unmarshal([]byte(defVal), &resp)
			//	return defVal, err
			//}
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		return handler(ctx, req)
	}
}

func (b *InterceptorBuilder) BuildClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		limited, err := b.limiter.Limit(ctx, b.key)

		if err != nil {
			return status.Errorf(codes.ResourceExhausted, "触发限流")
		}

		if limited {
			return status.Errorf(codes.ResourceExhausted, "触发限流")
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// 服务级别限流

func (b *InterceptorBuilder) BuildServerInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if strings.HasPrefix(info.FullMethod, "/UserService") {
			limited, err := b.limiter.Limit(ctx, "limiter:service:user:UserService")
			if err != nil {
				// err 不为nil，你要考虑你用保守的，还是用激进的策略
				// 这是保守的策略
				b.l.Error("判定限流出现问题")
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")

				// 这是激进的策略
				// return handler(ctx, req)
			}
			if limited {
				//defVal, ok := b.defaultValueMap[info.FullMethod]
				//if ok {
				//	err = json.Unmarshal([]byte(defVal), &resp)
				//	return defVal, err
				//}
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
		}

		return handler(ctx, req)
	}
}
