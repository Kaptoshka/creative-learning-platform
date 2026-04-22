package grpc

import (
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mapGRPCError converts gRPC status errors to domain sentinel errors.
func mapGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return domain.ErrInternal
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return domain.ErrUnauthenticated
	case codes.PermissionDenied:
		return domain.ErrPermissionDenied
	case codes.NotFound:
		return domain.ErrNotFound
	case codes.AlreadyExists:
		return domain.ErrAlreadyExists
	case codes.InvalidArgument:
		return domain.ErrInvalidArgument
	default:
		return domain.ErrInternal
	}
}
