package domain

import "time"

// --- Enums ---

type SubmissionStatus int32

const (
	SubmissionStatusUnspecified SubmissionStatus = 0
	SubmissionStatusNotStarted  SubmissionStatus = 1
	SubmissionStatusInProgress  SubmissionStatus = 2
	SubmissionStatusSubmitted   SubmissionStatus = 3
	SubmissionStatusGraded      SubmissionStatus = 4
	SubmissionStatusReturned    SubmissionStatus = 5
)

// --- Entities ---

type AssignmentTemplate struct {
	ID           string `json:"id"`
	CreatorID    string `json:"creator_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	WidgetID     string `json:"widget_id"`
	WidgetConfig map[string]any `json:"widget_config"`
	DueDate      *time.Time `json:"due_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AssignmentTemplateLight struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	WidgetType string `json:"widget_type"`
	DueDate    *time.Time `json:"due_date"`
}

type AssignmentTarget struct {
	// only one of the following fields should be set
	GroupID   string `json:"group_id"`
	StudentID string `json:"student_id"`
}

type Submission struct {
	ID            string `json:"id"`
	TemplateID    string `json:"template_id"`
	StudentID     string `json:"student_id"`
	Status        SubmissionStatus `json:"status"`
	StartedAt     *time.Time `json:"started_at"`
	SubmittedAt   *time.Time `json:"submitted_at"`
	LatestVersion *SubmissionVersionLight `json:"latest_version"`
}

type SubmissionVersion struct {
	ID            string `json:"id"`
	VersionNumber int32 `json:"version_number"`
	Payload       map[string]any `json:"payload"`
	TimeSpent     time.Duration `json:"time_spent"`
	IsAutosave    bool `json:"is_autosave"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SubmissionVersionLight struct {
	ID            string `json:"id"`
	VersionNumber int32 `json:"version_number"`
	CreatedAt     time.Time `json:"created_at"`
}

type Feedback struct {
	ID          string `json:"id"`
	VersionID   string `json:"version_id"`
	GraderID    string `json:"grader_id"`
	TextContent string `json:"text_content"`
	Payload     map[string]any `json:"payload"`
	IsPublished bool `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// --- Teacher: template management ---

type CreateAssignmentRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	WidgetID     string `json:"widget_id"`
	WidgetConfig map[string]any `json:"widget_config"`
	DueDate      *time.Time `json:"due_date"`
	Targets      []AssignmentTarget `json:"targets"`
}

type CreateAssignmentResponse struct {
	ID string `json:"id"`
}

type UpdateAssignmentRequest struct {
	AssignmentID string `json:"assignment_id"`
	Template     AssignmentTemplate `json:"template"`
	UpdateMask   []string `json:"update_mask"`
	Targets      []AssignmentTarget `json:"targets"`
}

type UpdateAssignmentResponse struct {
	Template AssignmentTemplate `json:"template"`
}

type DeleteAssignmentRequest struct {
	ID string `json:"id"`
}

type GetAssignmentRequest struct {
	ID string `json:"id"`
}

type GetAssignmentResponse struct {
	Template AssignmentTemplate `json:"template"`
	Targets  []AssignmentTarget `json:"targets"`
}

// --- Teacher: submissions & feedback ---

type ListAssignmentsRequest struct {
	PageSize  int32 `json:"page_size"`
	PageToken string `json:"page_token"`
	CreatorID string `json:"creator_id"`
}

type ListAssignmentsResponse struct {
	Items         []AssignmentTemplateLight  `json:"items"`
	NextPageToken string `json:"next_page_token"`
}

type ListAssignmentSubmissionsRequest struct {
	TemplateID   string `json:"template_id"`
	PageSize     int32 `json:"page_size"`
	PageToken    string `json:"page_token"`
	StatusFilter SubmissionStatus `json:"status_filter"`
}

type ListAssignmentSubmissionsResponse struct {
	Items         []Submission `json:"items"`
	NextPageToken string `json:"next_page_token"`
}

type GetStudentSubmissionRequest struct {
	SubmissionID string `json:"submission_id"`
}

type GetStudentSubmissionResponse struct {
	Template   AssignmentTemplate `json:"template"`
	Submission Submission `json:"submission"`
	History    []SubmissionVersion `json:"history"`
	Feedback   []Feedback `json:"feedback"`
}

type ProvideFeedbackRequest struct {
	SubmissionID string `json:"submission_id"`
	VersionID    string `json:"version_id"`
	TextContent  string `json:"text_content"`
	Payload      map[string]any `json:"payload"`
	IsPublished  bool `json:"is_published"`
}

// --- Student: assignment workflow ---

type ListMyAssignmentsRequest struct {
	PageSize     int32 `json:"page_size"`
	PageToken    string `json:"page_token"`
	StatusFilter SubmissionStatus `json:"status_filter"`
}

type ListMyAssignmentsItem struct {
	Template    AssignmentTemplateLight `json:"template"`
	Status      SubmissionStatus `json:"status"`
	HasFeedback bool `json:"has_feedback"`
}

type ListMyAssignmentsResponse struct {
	Items         []ListMyAssignmentsItem `json:"items"`
	NextPageToken string `json:"next_page_token"`
}

type StartAssignmentRequest struct {
	TemplateID string `json:"template_id"`
}

type StartAssignmentResponse struct {
	SubmissionID string `json:"submission_id"`
	StartedAt    time.Time `json:"started_at"`
}

type SaveAssignmentDraftRequest struct {
	SubmissionID string `json:"submission_id"`
	Payload      map[string]any `json:"payload"`
	TimeSpent    time.Duration `json:"time_spent"`
}

type SaveAssignmentDraftResponse struct {
	VersionID string `json:"version_id"`
	SavedAt   time.Time `json:"saved_at"`
}

type SubmitAssignmentRequest struct {
	SubmissionID string `json:"submission_id"`
	Payload      map[string]any `json:"payload"`
	TimeSpent    time.Duration `json:"time_spent"`
}

type SubmitAssignmentResponse struct {
	VersionID string `json:"version_id"`
	Status    SubmissionStatus `json:"status"`
}
