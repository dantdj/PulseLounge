package httpapi

import (
	"net/http"
	"strings"

	"pulselounge/internal/channels"
	"pulselounge/internal/messages"
)

func NewRouter(uiHandler http.Handler, messageService messages.Service, channelService channels.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)

	channelsHandler := NewChannelsHandler(channelService)
	messagesHandler := NewMessagesHandler(messageService)
	mux.HandleFunc("/api/channels", channelsHandler.Channels)
	mux.HandleFunc("/api/channels/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/messages") {
			messagesHandler.ChannelMessages(w, r)
			return
		}
		channelsHandler.Channel(w, r)
	})

	mux.HandleFunc("/api/messages/", messagesHandler.Message)

	mux.Handle("/", uiHandler)
	return mux
}
