package incoming

import "net/http"

type SSOHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	LogoutAll(w http.ResponseWriter, r *http.Request)

	// Admin only
	RegisterApp(w http.ResponseWriter, r *http.Request)
	DeactivateApp(w http.ResponseWriter, r *http.Request)
}

type AssignmentHandler interface {
	// Teacher: template management
	CreateAssignment(w http.ResponseWriter, r *http.Request)
	UpdateAssignment(w http.ResponseWriter, r *http.Request)
	DeleteAssignment(w http.ResponseWriter, r *http.Request)
	GetAssignment(w http.ResponseWriter, r *http.Request)

	// Teacher: submissions & feedback
	ListAssignments(w http.ResponseWriter, r *http.Request)
	ListAssignmentSubmissions(w http.ResponseWriter, r *http.Request)
	GetStudentSubmission(w http.ResponseWriter, r *http.Request)
	ProvideFeedback(w http.ResponseWriter, r *http.Request)

	// Student: assignment workflow
	ListMyAssignments(w http.ResponseWriter, r *http.Request)
	StartAssignment(w http.ResponseWriter, r *http.Request)
	SaveAssignmentDraft(w http.ResponseWriter, r *http.Request)
	SubmitAssignment(w http.ResponseWriter, r *http.Request)
}
