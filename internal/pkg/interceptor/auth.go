package interceptor

import (
	"context"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(logger *log.Logger, jwtManager *jwt.JWTManager, authMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{
		logger:      logger,
		jwtManager:  jwtManager,
		authMethods: authMethods,
	}
}

type AuthInterceptor struct {
	logger      *log.Logger
	jwtManager  *jwt.JWTManager
	authMethods map[string]bool
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		i.logger.Info("--> unary interceptor: ", info.FullMethod)
		claims, err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		if claims != nil {
			ctx = context.WithValue(ctx, "ID", claims.ID)
		}
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (*jwt.UserClaims, error) {
	b, ok := i.authMethods[method]
	if !ok || !b {
		return nil, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	token := values[0]
	claims, err := i.jwtManager.Verify(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is invalid")
	}

	return claims, nil
}
