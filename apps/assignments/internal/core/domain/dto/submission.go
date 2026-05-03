package dto

import (
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/models"
)

type FullSubmission struct {
	Assignment *models.AssignmentTemplate
	Targets    []models.AssignmentTarget
	Submission *models.Submission
	Versions   []models.SubmissionVersion
	Feedbacks  []models.Feedback
}
