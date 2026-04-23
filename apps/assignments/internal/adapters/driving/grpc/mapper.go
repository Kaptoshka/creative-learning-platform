package grpc

import (
	"encoding/json"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain"
	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func convertSubmissionStatus(status domain.SubmissionStatus) tasksv1.SubmissionStatus {
	switch status {
	case domain.StatusNotStarted:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_NOT_STARTED
	case domain.StatusInProgress:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_IN_PROGRESS
	case domain.StatusSubmitted:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_SUBMITTED
	case domain.StatusGraded:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_GRADED
	case domain.StatusReturned:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_RETURNED
	default:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_UNSPECIFIED
	}
}

func convertProtoStatus(status tasksv1.SubmissionStatus) domain.SubmissionStatus {
	switch status {
	case tasksv1.SubmissionStatus_SUBMISSION_STATUS_NOT_STARTED:
		return domain.StatusNotStarted
	case tasksv1.SubmissionStatus_SUBMISSION_STATUS_IN_PROGRESS:
		return domain.StatusInProgress
	case tasksv1.SubmissionStatus_SUBMISSION_STATUS_SUBMITTED:
		return domain.StatusSubmitted
	case tasksv1.SubmissionStatus_SUBMISSION_STATUS_GRADED:
		return domain.StatusGraded
	case tasksv1.SubmissionStatus_SUBMISSION_STATUS_RETURNED:
		return domain.StatusReturned
	default:
		return domain.StatusNotSpecified
	}
}

func rawMessageToStructPB(raw json.RawMessage) (*structpb.Struct, error) {
	if len(raw) <= 0 {
		return &structpb.Struct{}, nil
	}
	s := &structpb.Struct{}
	if err := protojson.Unmarshal(raw, s); err != nil {
		return nil, fmt.Errorf("convert to structpb: %w", err)
	}
	return s, nil
}

func structPBToRawMessage(s *structpb.Struct) (json.RawMessage, error) {
	if s == nil {
		return json.RawMessage("{}"), nil
	}
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("convert to rawMessage: %w", err)
	}
	return b, nil
}
