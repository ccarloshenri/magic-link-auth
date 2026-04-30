package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"magic-link-auth/src/layers/main/processor"
)

type CreateMagicLinkController struct {
	processor *processor.CreateMagicLinkProcessor
}

func NewCreateMagicLinkController(p *processor.CreateMagicLinkProcessor) *CreateMagicLinkController {
	return &CreateMagicLinkController{processor: p}
}

func (c *CreateMagicLinkController) Handle(w http.ResponseWriter, r *http.Request) {
	var input processor.CreateMagicLinkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := c.processor.Process(input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("magic link sent to %s", input.Email)})
}
