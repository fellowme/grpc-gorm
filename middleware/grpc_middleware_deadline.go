package deadline

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"grpc-gorm/tools/logSetup"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ctx.Err() == context.Canceled {
			logSetup.MyLogger.Error("请求超时")
			return nil, errors.New("超时取消")
		}
		resp, err := handler(ctx, req)
		return resp, err
	}
}
