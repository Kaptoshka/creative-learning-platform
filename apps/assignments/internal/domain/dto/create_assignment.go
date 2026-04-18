package dto

import (
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"

	"github.com/google/uuid"
)

type CreateAssignment struct {
	Title        string       `db:"title"`
	Description  string       `db:"description"`
	WidgetID     uuid.UUID    `db:"widget_id"`
	WidgetConfig domain.JSONB `db:"widget_config"`
	DueDate      *time.Time   `db:"due_date"`
	Targets      []*Target    `db:"targets"`
}
