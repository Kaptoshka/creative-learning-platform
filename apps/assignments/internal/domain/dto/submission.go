package dto

import (
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/models"
)

type FullSubmission struct {
	Template   *models.AssignmentTemplate
	Targets    []*models.AssignmentTarget
	Submission *models.Submission
	Versions   []*models.SubmissionVersion
	Feedbacks  []*models.Feedback
}
