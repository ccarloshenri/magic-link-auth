package controller

import (
	"errors"
	"net/http"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/bo"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/processor"
)

type ValidateMagicLinkController struct {
	processor *processor.ValidateMagicLinkProcessor
}

func NewValidateMagicLinkController(p *processor.ValidateMagicLinkProcessor) *ValidateMagicLinkController {
	return &ValidateMagicLinkController{processor: p}
}

func (c *ValidateMagicLinkController) Handle(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	output, err := c.processor.Process(token)
	if err != nil {
		writeJSON(w, resolveStatus(err), map[string]string{"error": err.Error()})
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
