package models

import (
	"time"

	"github.com/google/uuid"
)

type AssignmentTemplate struct {
	ID           uuid.UUID  `db:"id"`
	CreatorID    uuid.UUID  `db:"creator_id"`
	Title        string     `db:"title"`
	Description  string     `db:"description"`
	WidgetID     uuid.UUID  `db:"widget_id"`
	WidgetConfig JSONB      `db:"widget_config"`
	DueDate      *time.Time `db:"due_date"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

type AssignmentTemplateLight struct {
	ID         uuid.UUID        `db:"id"`
	Title      string           `db:"title"`
	Status     SubmissionStatus `db:"status"`
	WidgetType string           `db:"widget_type"`
	DueDate    *time.Time       `db:"due_date"`
}

type AssignmentTarget struct {
	ID         uuid.UUID  `db:"id"`
	TemplateID uuid.UUID  `db:"template_id"`
	GroupID    *uuid.UUID `db:"group_id"`
	StudentID  *uuid.UUID `db:"student_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type TargetDTO struct {
	TemplateID uuid.UUID  `db:"template_id"`
	GroupID    *uuid.UUID `db:"group_id"`
	StudentID  *uuid.UUID `db:"student_id"`
}

type CreateAssignmentDTO struct {
	Title        string       `db:"title"`
	Description  string       `db:"description"`
	WidgetID     uuid.UUID    `db:"widget_id"`
	WidgetConfig JSONB        `db:"widget_config"`
	DueDate      *time.Time   `db:"due_date"`
	Targets      []*TargetDTO `db:"targets"`
}
