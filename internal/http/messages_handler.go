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

func (h MessagesHandler) Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		ID      int    `json:"id"`
		NewBody string `json:"newBody"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.Edit(r.Context(), req.ID, req.NewBody)
	if err != nil {
		switch {
		case errors.Is(err, messages.ErrEmptyBody):
			writeJSONError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, messages.ErrMessageNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "failed to edit message")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
