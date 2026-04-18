package grpc

import (
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
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
