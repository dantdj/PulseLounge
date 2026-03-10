package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"pulselounge/internal/messages"
)

type MessagesHandler struct {
	service messages.Service
}

func NewMessagesHandler(service messages.Service) MessagesHandler {
	return MessagesHandler{service: service}
}

func (h MessagesHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	result, err := h.service.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to query messages")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h MessagesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req.Body)
	if err != nil {
		if errors.Is(err, messages.ErrEmptyBody) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to create message")
		return
	}

	writeJSON(w, http.StatusCreated, result)
}
