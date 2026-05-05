package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pulselounge/internal/media"
	"pulselounge/internal/messages"
)

const devAuthorID int64 = 1

type MessagesHandler struct {
	service messages.Service
	store   media.Store
}

func NewMessagesHandler(service messages.Service, store media.Store) MessagesHandler {
	return MessagesHandler{service: service, store: store}
}

func (h MessagesHandler) ChannelMessages(w http.ResponseWriter, r *http.Request) {
	channelID, ok := channelIDFromPath(r.URL.Path)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid channel id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.List(w, r, channelID)
	case http.MethodPost:
		h.Create(w, r, channelID)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h MessagesHandler) Message(w http.ResponseWriter, r *http.Request) {
	messageID, ok := messageIDFromPath(r.URL.Path)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid message id")
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.Edit(w, r, messageID)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h MessagesHandler) List(w http.ResponseWriter, r *http.Request, channelID int64) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	result, err := h.service.ListByChannel(r.Context(), channelID)
	if err != nil {
		log.Printf("failed to query messages for channel %d: %v", channelID, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to query messages")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h MessagesHandler) Create(w http.ResponseWriter, r *http.Request, channelID int64) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Body     string `json:"body"`
		ImageKey string `json:"imageKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.CreateInChannel(r.Context(), channelID, devAuthorID, req.Body, req.ImageKey)
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

func (h MessagesHandler) Edit(w http.ResponseWriter, r *http.Request, messageID int64) {
	if r.Method != http.MethodPut {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		NewBody string `json:"newBody"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.Edit(r.Context(), messageID, req.NewBody)
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

func channelIDFromPath(path string) (int64, bool) {
	rest := strings.TrimPrefix(path, "/api/channels/")
	parts := strings.Split(rest, "/")
	if len(parts) != 2 || parts[1] != "messages" {
		return 0, false
	}
	return positiveInt64(parts[0])
}

func messageIDFromPath(path string) (int64, bool) {
	rest := strings.TrimPrefix(path, "/api/messages/")
	if rest == "" || strings.Contains(rest, "/") {
		return 0, false
	}
	return positiveInt64(rest)
}

func positiveInt64(value string) (int64, bool) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}
