package main

import (
	"database/sql"
	"io/fs"
	"net/http"
)

type app struct {
	db *sql.DB
}

func (a app) routes(uiFS fs.FS) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/messages", a.listMessagesHandler)
	mux.HandleFunc("/", spaHandler(uiFS))
	return mux
}
