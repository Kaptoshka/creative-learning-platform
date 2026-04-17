package submission

import (
	"log/slog"
)

type SubmissionService struct {
	log                *slog.Logger
	submissionSaver    SubmissionSaver
	submissionProvider SubmissionProvider
}

type SubmissionSaver interface {
}

type SubmissionProvider interface {
}

func New(
	log *slog.Logger,
	submissionSaver SubmissionSaver,
	submissionProvider SubmissionProvider,
) *SubmissionService {
	return &SubmissionService{
		log:                log,
		submissionSaver:    submissionSaver,
		submissionProvider: submissionProvider,
	}
}
