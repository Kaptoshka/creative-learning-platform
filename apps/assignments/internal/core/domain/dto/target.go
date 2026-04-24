package dto

import (
	"github.com/google/uuid"
)

type Target struct {
	TemplateID uuid.UUID  `db:"template_id"`
	GroupID    *uuid.UUID `db:"group_id"`
	StudentID  *uuid.UUID `db:"student_id"`
}
