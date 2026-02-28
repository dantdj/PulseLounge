package main

import (
	"database/sql"
	"net/http"
)

type app struct {
	db *sql.DB
}

func (a app) routes(uiHandler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/messages", a.listMessagesHandler)
	mux.Handle("/", uiHandler)
	return mux
}
