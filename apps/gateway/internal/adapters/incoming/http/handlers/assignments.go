package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/middleware"
	httputil "github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/utils"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/outgoing"
)

type AssignmentsHandler struct {
	uc  outgoing.AssignmentService
	log *slog.Logger
}

func NewAssignmentsHandler(uc outgoing.AssignmentService, log *slog.Logger) *AssignmentsHandler {
	return &AssignmentsHandler{
		uc:  uc,
		log: log,
	}
}

// --- Teacher: template management ---

// CreateAssignment godoc
// @Summary Create assignment template
// @Description Creates a new assignment template for teachers
// @Tags Assignments
// @Accept json
// @Produce json
// @Param request body domain.CreateAssignmentRequest true "Assignment template details"
// @Success 201 {object} domain.CreateAssignmentResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Router /assignments [post].
func (h *AssignmentsHandler) CreateAssignment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req domain.CreateAssignmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.CreateAssignment(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.Created(h.log, w, res)
}

// UpdateAssignment godoc
// @Summary Update assignment template
// @Description Updates an existing assignment template with optional field mask
// @Tags Assignments
// @Accept json
// @Produce json
// @Param request body domain.UpdateAssignmentRequest true "Assignment update details"
// @Success 200 {object} domain.UpdateAssignmentResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /assignments [patch].
func (h *AssignmentsHandler) UpdateAssignment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req domain.UpdateAssignmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.UpdateAssignment(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// DeleteAssignment godoc
// @Summary Delete assignment template
// @Description Deletes an assignment template by ID
// @Tags Assignments
// @Param id path string true "Assignment Template ID"
// @Success 204 "No Content"
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /assignments/{id} [delete].
func (h *AssignmentsHandler) DeleteAssignment(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	if err := h.uc.DeleteAssignment(
		r.Context(), domain.DeleteAssignmentRequest{ID: id},
	); err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.NoContent(w)
}

// GetAssignment godoc
// @Summary Get assignment template
// @Description Retrieves a single assignment template by ID
// @Tags Assignments
// @Produce json
// @Param id path string true "Assignment Template ID"
// @Success 200 {object} domain.GetAssignmentResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /assignments/{id} [get].
func (h *AssignmentsHandler) GetAssignment(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.GetAssignment(
		r.Context(),
		domain.GetAssignmentRequest{ID: id},
	)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// --- Teacher: submissions & feedback ---

// ListAssignments godoc
// @Summary List assignment templates
// @Description Returns a paginated list of assignment templates, optionally filtered by creator
// @Tags Assignments
// @Produce json
// @Param page_size query int false "Page size (10, 25, 50, default 20)"
// @Param page_token query string false "Pagination token"
// @Param creator_id query string false "Filter by creator ID"
// @Success 200 {object} domain.ListAssignmentsResponse
// @Failure 401 {object} domain.Error
// @Router /assignments [get].
func (h *AssignmentsHandler) ListAssignments(
	w http.ResponseWriter,
	r *http.Request,
) {
	res, err := h.uc.ListAssignments(r.Context(), domain.ListAssignmentsRequest{
		PageSize:  pageSize(r),
		PageToken: r.URL.Query().Get("page_token"),
		CreatorID: r.URL.Query().Get("creator_id"),
	})
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// ListAssignmentSubmissions godoc
// @Summary List assignment submissions
// @Description Returns a paginated list of student submissions for a specific assignment template
// @Tags Assignments
// @Produce json
// @Param template_id path string true "Assignment Template ID"
// @Param page_size query int false "Page size (10, 25, 50, default 20)"
// @Param page_token query string false "Pagination token"
// @Success 200 {object} domain.ListAssignmentSubmissionsResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Router /assignments/{template_id}/submissions [get].
func (h *AssignmentsHandler) ListAssignmentSubmissions(
	w http.ResponseWriter,
	r *http.Request,
) {
	templateID := r.PathValue("template_id")
	if templateID == "" {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.ListAssignmentSubmissions(
		r.Context(),
		domain.ListAssignmentSubmissionsRequest{
			TemplateID: templateID,
			PageSize:   pageSize(r),
			PageToken:  r.URL.Query().Get("page_token"),
		},
	)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// GetStudentSubmission godoc
// @Summary Get student submission
// @Description Retrieves a student's submission with full history and feedback
// @Tags Assignments
// @Produce json
// @Param submission_id path string true "Submission ID"
// @Success 200 {object} domain.GetStudentSubmissionResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /assignments/submissions/{submission_id} [get].
func (h *AssignmentsHandler) GetStudentSubmission(
	w http.ResponseWriter,
	r *http.Request,
) {
	submissionID := r.PathValue("submission_id")
	if submissionID == "" {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.GetStudentSubmission(
		r.Context(),
		domain.GetStudentSubmissionRequest{
			SubmissionID: submissionID,
		},
	)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// ProvideFeedback godoc
// @Summary Provide feedback on submission
// @Description Teacher provides feedback (grade, comments) for a student's submission version
// @Tags Assignments
// @Accept json
// @Param submission_id path string true "Submission ID"
// @Param request body domain.ProvideFeedbackRequest true "Feedback details"
// @Success 204 "No Content"
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /assignments/submissions/{submission_id}/feedback [post].
func (h *AssignmentsHandler) ProvideFeedback(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req domain.ProvideFeedbackRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	req.SubmissionID = r.PathValue("submission_id")

	if err := h.uc.ProvideFeedback(r.Context(), req); err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.NoContent(w)
}

// --- Student: assignment workflow ---

// ListMyAssignments godoc
// @Summary List my assignments
// @Description Returns a paginated list of assignments assigned to the current student
// @Tags Assignments
// @Produce json
// @Param page_size query int false "Page size (10, 25, 50, default 20)"
// @Param page_token query string false "Pagination token"
// @Success 200 {object} domain.ListMyAssignmentsResponse
// @Failure 401 {object} domain.Error
// @Router /students/me/assignments [get].
func (h *AssignmentsHandler) ListMyAssignments(
	w http.ResponseWriter,
	r *http.Request,
) {
	res, err := h.uc.ListMyAssignments(
		r.Context(),
		domain.ListMyAssignmentsRequest{
			PageSize:  pageSize(r),
			PageToken: r.URL.Query().Get("page_token"),
		},
	)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// StartAssignment godoc
// @Summary Start assignment
// @Description Student starts working on an assignment template, creates a submission
// @Tags Assignments
// @Produce json
// @Param template_id path string true "Assignment Template ID"
// @Success 201 {object} domain.StartAssignmentResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /students/me/assignments/{template_id} [post].
func (h *AssignmentsHandler) StartAssignment(
	w http.ResponseWriter,
	r *http.Request,
) {
	templateID := r.PathValue("template_id")
	if templateID == "" {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	_ = middleware.TokenFromCtx(r.Context())

	res, err := h.uc.StartAssignment(r.Context(), domain.StartAssignmentRequest{
		TemplateID: templateID,
	})
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.Created(h.log, w, res)
}

// SaveAssignmentDraft godoc
// @Summary Save assignment draft
// @Description Student saves draft progress for a submission (autosave)
// @Tags Assignments
// @Accept json
// @Produce json
// @Param submission_id path string true "Submission ID"
// @Param request body domain.SaveAssignmentDraftRequest true "Draft content"
// @Success 200 {object} domain.SaveAssignmentDraftResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /students/me/assignments/submissions/{submission_id}/draft [put].
func (h *AssignmentsHandler) SaveAssignmentDraft(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req domain.SaveAssignmentDraftRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	req.SubmissionID = r.PathValue("submission_id")

	res, err := h.uc.SaveAssignmentDraft(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// SubmitAssignment godoc
// @Summary Submit assignment
// @Description Student submits final version of their assignment for grading
// @Tags Assignments
// @Accept json
// @Produce json
// @Param submission_id path string true "Submission ID"
// @Param request body domain.SubmitAssignmentRequest true "Final submission content"
// @Success 200 {object} domain.SubmitAssignmentResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Router /students/me/assignments/submissions/{submission_id}/submit [post].
func (h *AssignmentsHandler) SubmitAssignment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req domain.SubmitAssignmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	req.SubmissionID = r.PathValue("submission_id")

	res, err := h.uc.SubmitAssignment(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// --- Helpers ---

const (
	tenPageSize        = 10
	twentyPageSize     = 20
	twentyFivePageSize = 25
	fiftyPageSize      = 50
)

func pageSize(r *http.Request) int32 {
	switch r.URL.Query().Get("page_size") {
	case "10":
		return tenPageSize
	case "25":
		return twentyFivePageSize
	case "50":
		return fiftyPageSize
	default:
		return twentyPageSize
	}
}
