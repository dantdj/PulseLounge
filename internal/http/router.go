package httpapi

import (
	"net/http"

	"pulselounge/internal/messages"
)

func NewRouter(uiHandler http.Handler, messageService messages.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)

	messagesHandler := NewMessagesHandler(messageService)
	mux.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			messagesHandler.List(w, r)
		case http.MethodPost:
			messagesHandler.Create(w, r)
		case http.MethodPut:
			messagesHandler.Edit(w, r)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.Handle("/", uiHandler)
	return mux
}
