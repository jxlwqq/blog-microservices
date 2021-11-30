package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

func NewInterceptor(jwtManager *JWTManager, methods map[string]bool) *Interceptor {
	return &Interceptor{
		jwtManager: jwtManager,
		methods:    methods,
	}
}

type Interceptor struct {
	jwtManager *JWTManager
	methods    map[string]bool
}

func (i Interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Println("--> unary interceptor: ", info.FullMethod)
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

func (i *Interceptor) authorize(ctx context.Context, method string) (*UserClaims, error) {
	b, ok := i.methods[method]
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
