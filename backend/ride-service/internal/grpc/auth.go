package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const internalTokenHeader = "x-internal-service-token"

func InternalAuthInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if token == "" {
			return handler(ctx, request)
		}

		metadata, ok := metadata.FromIncomingContext(ctx)
		if !ok || !hasToken(metadata.Get(internalTokenHeader), token) {
			return nil, status.Error(codes.Unauthenticated, "missing internal service token")
		}

		return handler(ctx, request)
	}
}

func hasToken(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}

	return false
}
