package models

import (
	"time"
)

type Assignment struct {
	ID        int64
	TeacherID int64                  `json:"teacher_id"`
	Title     string                 `json:"title"`
	Content   map[string]interface{} `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Deadline  time.Time              `json:"deadline"`
}
