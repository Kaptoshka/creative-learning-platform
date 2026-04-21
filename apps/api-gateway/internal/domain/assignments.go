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
	ID          string
	CreatorID   string
	Title       string
	Description string
	WidgetID    string
	WidgetConfig map[string]any
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AssignmentTemplateLight struct {
	ID         string
	Title      string
	WidgetType string
	DueDate    *time.Time
}

type AssignmentTarget struct {
	// only one of the following fields should be set
	GroupID   string
	StudentID string
}

type Submission struct {
	ID            string
	TemplateID    string
	StudentID     string
	Status        SubmissionStatus
	StartedAt     *time.Time
	SubmittedAt   *time.Time
	LatestVersion *SubmissionVersion
}

type SubmissionVersion struct {
	ID             string
	VersionNumber  int32
	Payload        map[string]any
	TimeSpent      time.Duration
	IsAutosave     bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Feedback struct {
	ID          string
	VersionID   string
	GraderID    string
	TextContent string
	Payload     map[string]any
	IsPublished bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// --- Teacher: template management ---

type CreateAssignmentRequest struct {
	Title        string
	Description  string
	WidgetID     string
	WidgetConfig map[string]any
	DueDate      *time.Time
	Targets      []AssignmentTarget
}

type CreateAssignmentResponse struct {
	ID string
}

type UpdateAssignmentRequest struct {
	AssignmentID string
	Template     AssignmentTemplate
	UpdateMask   []string
	Targets      []AssignmentTarget
}

type UpdateAssignmentResponse struct {
	Template AssignmentTemplate
}

type DeleteAssignmentRequest struct {
	ID string
}

type GetAssignmentRequest struct {
	ID string
}

type GetAssignmentResponse struct {
	Template AssignmentTemplate
	Targets  []AssignmentTarget
}

// --- Teacher: submissions & feedback ---

type ListAssignmentsRequest struct {
	PageSize  int32
	PageToken string
	CreatorID string
}

type ListAssignmentsResponse struct {
	Items         []AssignmentTemplateLight
	NextPageToken string
}

type ListAssignmentSubmissionsRequest struct {
	TemplateID   string
	PageSize     int32
	PageToken    string
	StatusFilter SubmissionStatus
}

type ListAssignmentSubmissionsResponse struct {
	Items         []Submission
	NextPageToken string
}

type GetStudentSubmissionRequest struct {
	SubmissionID string
}

type GetStudentSubmissionResponse struct {
	Template  AssignmentTemplate
	Submission Submission
	History   []SubmissionVersion
	Feedback  []Feedback
}

type ProvideFeedbackRequest struct {
	SubmissionID string
	VersionID    string
	TextContent  string
	Payload      map[string]any
	IsPublished  bool
}

// --- Student: assignment workflow ---

type ListMyAssignmentsRequest struct {
	PageSize     int32
	PageToken    string
	StatusFilter SubmissionStatus
}

type ListMyAssignmentsItem struct {
	Template    AssignmentTemplateLight
	Status      SubmissionStatus
	HasFeedback bool
}

type ListMyAssignmentsResponse struct {
	Items         []ListMyAssignmentsItem
	NextPageToken string
}

type StartAssignmentRequest struct {
	TemplateID string
}

type StartAssignmentResponse struct {
	SubmissionID string
	StartedAt    time.Time
}

type SaveAssignmentDraftRequest struct {
	SubmissionID string
	Payload      map[string]any
	TimeSpent    time.Duration
}

type SaveAssignmentDraftResponse struct {
	VersionID string
	SavedAt   time.Time
}

type SubmitAssignmentRequest struct {
	SubmissionID string
	Payload      map[string]any
	TimeSpent    time.Duration
}

type SubmitAssignmentResponse struct {
	VersionID string
	Status    SubmissionStatus
}
