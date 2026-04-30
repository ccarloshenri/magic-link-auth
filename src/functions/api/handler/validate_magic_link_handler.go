package handler

import (
	"errors"
	"net/http"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/bo"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/processor"
)

type ValidateMagicLinkHandler struct {
	processor *processor.ValidateMagicLinkProcessor
}

func NewValidateMagicLinkHandler(p *processor.ValidateMagicLinkProcessor) *ValidateMagicLinkHandler {
	return &ValidateMagicLinkHandler{processor: p}
}

func (h *ValidateMagicLinkHandler) Handle(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	output, err := h.processor.Process(token)
	if err != nil {
		status := resolveStatus(err)
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, output)
}

func resolveStatus(err error) int {
	switch {
	case errors.Is(err, bo.ErrTokenNotFound):
		return http.StatusNotFound
	case errors.Is(err, bo.ErrTokenExpired), errors.Is(err, bo.ErrTokenUsed):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusBadRequest
	}
}
