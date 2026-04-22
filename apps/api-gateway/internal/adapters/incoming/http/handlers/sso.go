package handlers

import (
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/httputil"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/outgoing"
)

type SSOHandler struct {
	uc outgoing.SSOService
}

func NewSSOHandler(uc outgoing.SSOService) *SSOHandler {
	return &SSOHandler{uc: uc}
}

func (h *SSOHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.Register(r.Context(), req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.Created(w, res)
}

func (h *SSOHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.Login(r.Context(), req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}

func (h *SSOHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req domain.LogoutRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, domain.ErrInvalidArgument)
		return
	}

	res, err := h.uc.Logout(r.Context(), req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.OK(w, res)
}
