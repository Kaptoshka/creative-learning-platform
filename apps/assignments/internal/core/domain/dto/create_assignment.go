package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateAssignment struct {
	Title        string          `db:"title"`
	Description  string          `db:"description"`
	WidgetID     uuid.UUID       `db:"widget_id"`
	WidgetConfig json.RawMessage `db:"widget_config"`
	DueDate      *time.Time      `db:"due_date"`
	Targets      []*Target       `db:"targets"`
}
