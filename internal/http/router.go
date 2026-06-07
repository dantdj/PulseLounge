package httpapi

import (
	"log/slog"
	"net/http"
	"strings"

	"pulselounge/internal/channels"
	"pulselounge/internal/media"
	"pulselounge/internal/messages"
)

func NewRouter(uiHandler http.Handler, messageService messages.Service, channelService channels.Service, blobStore media.Store) http.Handler {
	return NewRouterWithLogger(uiHandler, messageService, channelService, blobStore, slog.Default())
}

func NewRouterWithLogger(uiHandler http.Handler, messageService messages.Service, channelService channels.Service, blobStore media.Store, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)

	channelsHandler := NewChannelsHandler(channelService)
	messagesHandler := NewMessagesHandler(messageService, blobStore)
	uploadHandler := NewUploadHandler(blobStore)

	mux.HandleFunc("/api/channels", channelsHandler.Channels)
	mux.HandleFunc("/api/channels/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/messages") {
			messagesHandler.ChannelMessages(w, r)
			return
		}
		channelsHandler.Channel(w, r)
	})

	mux.HandleFunc("/api/messages/", messagesHandler.Message)

	mux.HandleFunc("/api/upload", uploadHandler.Upload)

	mux.Handle("/", uiHandler)
	return withRequestLogging(logger, withPanicRecovery(logger, mux))
}
