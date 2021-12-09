package interceptor

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stonecutter/blog-microservices/internal/pkg/log"
	"github.com/stonecutter/blog-microservices/internal/pkg/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
	"time"
)

type MetricsInterceptor struct {
	m      metrics.Metrics
	logger *log.Logger
}

func NewMetricsInterceptor(logger *log.Logger, m metrics.Metrics) *MetricsInterceptor {
	return &MetricsInterceptor{
		m:      m,
		logger: logger,
	}
}

func (i MetricsInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		i.logger.Info("------> Metrics Unary Interceptor")
		start := time.Now()
		resp, err := handler(ctx, req)
		status := codes.OK
		if err != nil {
			status = ErrToGRPCCode(err)
		}
		i.m.ObserveResponseTime(int(status), info.FullMethod, info.FullMethod, time.Since(start).Seconds())
		i.m.IncHits(int(status), info.FullMethod, info.FullMethod)

		return resp, err
	}
}

func ErrToGRPCCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case strings.Contains(err.Error(), "Validate"):
		return codes.InvalidArgument
	case strings.Contains(err.Error(), "redis"):
		return codes.NotFound
	}
	return codes.Internal
}
