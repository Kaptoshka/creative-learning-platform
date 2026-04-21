package handlers

import (
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/httputil"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/middleware"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/outgoing"
)

type AssignmentsHandler struct {
	uc outgoing.AssignmentsService
}

func NewAssignmentsHandler(uc outgoing.AssignmentsService) *AssignmentsHandler {
	return &AssignmentsHandler{uc: uc}
}

// --- Teacher: template management ---

func (h *AssignmentsHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateAssignmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.CreateAssignment(r.Context(), req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.Created(w, res)
}

func (h *AssignmentsHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateAssignmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.UpdateAssignment(r.Context(), req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

func (h *AssignmentsHandler) DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	if err := h.uc.DeleteAssignment(r.Context(), domain.DeleteAssignmentRequest{ID: id}); err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.NoContent(w)
}

func (h *AssignmentsHandler) GetAssignment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.GetAssignment(r.Context(), domain.GetAssignmentRequest{ID: id})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

// --- Teacher: submissions & feedback ---

func (h *AssignmentsHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.ListAssignments(r.Context(), domain.ListAssignmentsRequest{
		PageSize:  pageSize(r),
		PageToken: r.URL.Query().Get("page_token"),
		CreatorID: r.URL.Query().Get("creator_id"),
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

func (h *AssignmentsHandler) ListAssignmentSubmissions(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("template_id")
	if templateID == "" {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.ListAssignmentSubmissions(r.Context(), domain.ListAssignmentSubmissionsRequest{
		TemplateID: templateID,
		PageSize:   pageSize(r),
		PageToken:  r.URL.Query().Get("page_token"),
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

func (h *AssignmentsHandler) GetStudentSubmission(w http.ResponseWriter, r *http.Request) {
	submissionID := r.PathValue("submission_id")
	if submissionID == "" {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.GetStudentSubmission(r.Context(), domain.GetStudentSubmissionRequest{
		SubmissionID: submissionID,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

func (h *AssignmentsHandler) ProvideFeedback(w http.ResponseWriter, r *http.Request) {
	var req domain.ProvideFeedbackRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	// submission_id берём из пути, не из тела
	req.SubmissionID = r.PathValue("submission_id")

	if err := h.uc.ProvideFeedback(r.Context(), req); err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.NoContent(w)
}

// --- Student: assignment workflow ---

func (h *AssignmentsHandler) ListMyAssignments(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.ListMyAssignments(r.Context(), domain.ListMyAssignmentsRequest{
		PageSize:  pageSize(r),
		PageToken: r.URL.Query().Get("page_token"),
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

func (h *AssignmentsHandler) StartAssignment(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("template_id")
	if templateID == "" {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	_ = middleware.TokenFromCtx(r.Context())

	res, err := h.uc.StartAssignment(r.Context(), domain.StartAssignmentRequest{
		TemplateID: templateID,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.Created(w, res)
}

func (h *AssignmentsHandler) SaveAssignmentDraft(w http.ResponseWriter, r *http.Request) {
	var req domain.SaveAssignmentDraftRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	req.SubmissionID = r.PathValue("submission_id")

	res, err := h.uc.SaveAssignmentDraft(r.Context(), req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

func (h *AssignmentsHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	var req domain.SubmitAssignmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	req.SubmissionID = r.PathValue("submission_id")

	res, err := h.uc.SubmitAssignment(r.Context(), req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

// --- Helpers ---

func pageSize(r *http.Request) int32 {
	switch r.URL.Query().Get("page_size") {
	case "10":
		return 10
	case "25":
		return 25
	case "50":
		return 50
	default:
		return 20
	}
}
