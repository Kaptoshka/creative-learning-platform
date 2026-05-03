package handlers

import (
	"log/slog"
	"net/http"

	httputil "github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/utils"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/outgoing"
)

type SSOHandler struct {
	uc  outgoing.SSOService
	log *slog.Logger
}

func NewSSOHandler(uc outgoing.SSOService, log *slog.Logger) *SSOHandler {
	return &SSOHandler{
		uc:  uc,
		log: log,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with email, password and personal information
// @Tags SSO
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration details"
// @Success 201 {object} domain.RegisterResponse
// @Failure 400 {object} domain.Error
// @Failure 409 {object} domain.Error
// @Router /sso/register [post].
func (h *SSOHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.Register(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.Created(h.log, w, res)
}

// Login godoc
// @Summary Login user
// @Description Authenticates user with email and password, returns access token
// @Tags SSO
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Router /sso/login [post].
func (h *SSOHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.Login(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

func (h *SSOHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.Refresh(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidates the user session token
// @Tags SSO
// @Accept json
// @Produce json
// @Param request body domain.LogoutRequest true "Logout request with token"
// @Success 200 {object} domain.LogoutResponse
// @Failure 400 {object} domain.Error
// @Router /sso/logout [post].
func (h *SSOHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req domain.LogoutRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.Logout(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.OK(h.log, w, res)
}

func (h *SSOHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	var req domain.LogoutAllRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	_, err := h.uc.LogoutAll(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.NoContent(w)
}

// --- Admin: App Management ---

func (h *SSOHandler) RegisterApp(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterAppRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.RegisterApp(r.Context(), req)
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.Created(h.log, w, res)
}

func (h *SSOHandler) DeactivateApp(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("app_id")
	if appID == "" {
		httputil.Error(h.log, w, domain.ErrInvalidArgument)
		return
	}

	_, err := h.uc.DeactivateApp(r.Context(), domain.DeactivateAppRequest{AppID: appID})
	if err != nil {
		httputil.Error(h.log, w, err)
		return
	}

	httputil.NoContent(w)
}
