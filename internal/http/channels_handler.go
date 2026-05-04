package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"pulselounge/internal/channels"
)

type ChannelsHandler struct {
	service channels.Service
}

func NewChannelsHandler(service channels.Service) ChannelsHandler {
	return ChannelsHandler{service: service}
}

func (h ChannelsHandler) Channels(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/channels" {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h ChannelsHandler) Channel(w http.ResponseWriter, r *http.Request) {
	channelID, ok := channelIDOnlyFromPath(r.URL.Path)
	if !ok {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.Delete(w, r, channelID)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h ChannelsHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to query channels")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h ChannelsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, channels.ErrEmptyName) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to create channel")
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h ChannelsHandler) Delete(w http.ResponseWriter, r *http.Request, channelID int64) {
	err := h.service.Delete(r.Context(), channelID)
	if err != nil {
		if errors.Is(err, channels.ErrChannelNotFound) {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to delete channel")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func channelIDOnlyFromPath(path string) (int64, bool) {
	rest := strings.TrimPrefix(path, "/api/channels/")
	if rest == "" || strings.Contains(rest, "/") {
		return 0, false
	}

	id, err := strconv.ParseInt(rest, 10, 64)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}
