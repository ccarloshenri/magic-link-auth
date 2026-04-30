package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/processor"
)

type CreateMagicLinkHandler struct {
	processor *processor.CreateMagicLinkProcessor
}

func NewCreateMagicLinkHandler(p *processor.CreateMagicLinkProcessor) *CreateMagicLinkHandler {
	return &CreateMagicLinkHandler{processor: p}
}

func (h *CreateMagicLinkHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var input processor.CreateMagicLinkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.processor.Process(input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("magic link sent to %s", input.Email)})
}
