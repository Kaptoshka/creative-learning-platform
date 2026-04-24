package router

import (
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/middleware"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/incoming"

	httpSwagger "github.com/swaggo/http-swagger"
)

func New(
	sso incoming.SSOHandler,
	assignments incoming.AssignmentHandler,
	mw *middleware.Middleware,
) http.Handler {
	mux := http.NewServeMux()

	// --- SSO: public routes ---
	mux.Handle(
		"POST /api/v1/auth/register",
		mw.Public(http.HandlerFunc(sso.Register)),
	)
	mux.Handle(
		"POST /api/v1/auth/login",
		mw.Public(http.HandlerFunc(sso.Login)),
	)
	mux.Handle(
		"POST /api/v1/auth/logout",
		mw.Public(http.HandlerFunc(sso.Logout)),
	)

	// --- Assignments: Teacher — template management ---
	mux.Handle(
		"GET /api/v1/assignments",
		mw.Protected(http.HandlerFunc(assignments.ListAssignments)),
	)
	mux.Handle(
		"POST /api/v1/assignments",
		mw.Protected(http.HandlerFunc(assignments.CreateAssignment)),
	)
	mux.Handle(
		"GET /api/v1/assignments/{id}",
		mw.Protected(http.HandlerFunc(assignments.GetAssignment)),
	)
	mux.Handle(
		"PATCH /api/v1/assignments/{id}",
		mw.Protected(http.HandlerFunc(assignments.UpdateAssignment)),
	)
	mux.Handle(
		"DELETE /api/v1/assignments/{id}",
		mw.Protected(http.HandlerFunc(assignments.DeleteAssignment)),
	)

	// --- Assignments: Teacher — submissions & feedback ---
	mux.Handle(
		"GET /api/v1/assignments/{template_id}/submissions",
		mw.Protected(http.HandlerFunc(assignments.ListAssignmentSubmissions)),
	)
	mux.Handle(
		"GET /api/v1/submissions/{submission_id}",
		mw.Protected(http.HandlerFunc(assignments.GetStudentSubmission)),
	)
	mux.Handle(
		"POST /api/v1/submissions/{submission_id}/feedback",
		mw.Protected(http.HandlerFunc(assignments.ProvideFeedback)),
	)

	// --- Assignments: Student — assignment workflow ---
	mux.Handle(
		"GET /api/v1/my/assignments",
		mw.Protected(http.HandlerFunc(assignments.ListMyAssignments)),
	)
	mux.Handle(
		"POST /api/v1/my/assignments/{template_id}/start",
		mw.Protected(http.HandlerFunc(assignments.StartAssignment)),
	)
	mux.Handle(
		"PUT /api/v1/my/submissions/{submission_id}/draft",
		mw.Protected(http.HandlerFunc(assignments.SaveAssignmentDraft)),
	)
	mux.Handle(
		"POST /api/v1/my/submissions/{submission_id}/submit",
		mw.Protected(http.HandlerFunc(assignments.SubmitAssignment)),
	)

	// --- Swagger UI ---
	mux.Handle(
		"/api/v1/swagger/",
		mw.Protected(httpSwagger.WrapHandler),
	)

	return mux
}
