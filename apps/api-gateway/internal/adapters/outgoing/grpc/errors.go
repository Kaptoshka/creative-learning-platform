package grpc

import (
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"

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
	case codes.FailedPrecondition:
		msg := st.Message()
		switch msg {
		case "submission is not in progress":
			return domain.ErrSubmissionClosed
		case "submission is not submitted":
			return domain.ErrSubmissionNotSubmitted
		}
		return domain.ErrInvalidArgument
	case codes.OutOfRange:
		return domain.ErrInvalidPageToken
	default:
		return domain.ErrInternal
	}
}
