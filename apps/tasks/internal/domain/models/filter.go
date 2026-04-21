package models

type Filter struct {
	Limit           int   `json:"limit"`
	Offset          int   `json:"offset"`
	TargetStudentID int64 `json:"target_student_id"`
	TeacherID       int64 `json:"teacher_id"`
}
