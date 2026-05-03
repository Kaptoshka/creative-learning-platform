package router

import (
	"log/slog"
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/middleware"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/incoming"

	httpSwagger "github.com/swaggo/http-swagger"
)

func New(
	log *slog.Logger,
	sso incoming.SSOHandler,
	assignments incoming.AssignmentHandler,
	mw *middleware.Middleware,
) http.Handler {
	mux := http.NewServeMux()

	// --- SSO: public routes ---
	mux.Handle(
		"POST /api/v1/auth/register",
		mw.Public(log, http.HandlerFunc(sso.Register)),
	)
	mux.Handle(
		"POST /api/v1/auth/login",
		mw.Public(log, http.HandlerFunc(sso.Login)),
	)
	mux.Handle(
		"POST /api/v1/auth/refresh",
		mw.Public(log, http.HandlerFunc(sso.Refresh)),
	)

	// --- SSO: need token ---
	mux.Handle(
		"POST /api/v1/auth/logout",
		mw.Protected(log, http.HandlerFunc(sso.Logout)),
	)
	mux.Handle(
		"POST /api/v1/auth/logout-all",
		mw.Protected(log, http.HandlerFunc(sso.LogoutAll)),
	)

	// --- SSO: admin only ---
	mux.Handle(
		"POST /api/v1/admin/apps",
		mw.Protected(log, http.HandlerFunc(sso.RegisterApp)),
	)
	mux.Handle(
		"POST /api/v1/admin/apps/{app_id}/deactivate",
		mw.Protected(log, http.HandlerFunc(sso.DeactivateApp)),
	)

	// --- Assignments: Teacher — template management ---
	mux.Handle(
		"GET /api/v1/assignments",
		mw.Protected(log, http.HandlerFunc(assignments.ListAssignments)),
	)
	mux.Handle(
		"POST /api/v1/assignments",
		mw.Protected(log, http.HandlerFunc(assignments.CreateAssignment)),
	)
	mux.Handle(
		"GET /api/v1/assignments/{id}",
		mw.Protected(log, http.HandlerFunc(assignments.GetAssignment)),
	)
	mux.Handle(
		"PATCH /api/v1/assignments/{id}",
		mw.Protected(log, http.HandlerFunc(assignments.UpdateAssignment)),
	)
	mux.Handle(
		"DELETE /api/v1/assignments/{id}",
		mw.Protected(log, http.HandlerFunc(assignments.DeleteAssignment)),
	)

	// --- Assignments: Teacher — submissions & feedback ---
	mux.Handle(
		"GET /api/v1/assignments/{template_id}/submissions",
		mw.Protected(log, http.HandlerFunc(assignments.ListAssignmentSubmissions)),
	)
	mux.Handle(
		"GET /api/v1/submissions/{submission_id}",
		mw.Protected(log, http.HandlerFunc(assignments.GetStudentSubmission)),
	)
	mux.Handle(
		"POST /api/v1/submissions/{submission_id}/feedback",
		mw.Protected(log, http.HandlerFunc(assignments.ProvideFeedback)),
	)

	// --- Assignments: Student — assignment workflow ---
	mux.Handle(
		"GET /api/v1/my/assignments",
		mw.Protected(log, http.HandlerFunc(assignments.ListMyAssignments)),
	)
	mux.Handle(
		"POST /api/v1/my/assignments/{template_id}/start",
		mw.Protected(log, http.HandlerFunc(assignments.StartAssignment)),
	)
	mux.Handle(
		"PUT /api/v1/my/submissions/{submission_id}/draft",
		mw.Protected(log, http.HandlerFunc(assignments.SaveAssignmentDraft)),
	)
	mux.Handle(
		"POST /api/v1/my/submissions/{submission_id}/submit",
		mw.Protected(log, http.HandlerFunc(assignments.SubmitAssignment)),
	)

	// --- Swagger UI ---
	mux.Handle(
		"/api/v1/swagger/",
		mw.Protected(log, httpSwagger.WrapHandler),
	)

	return mux
}
