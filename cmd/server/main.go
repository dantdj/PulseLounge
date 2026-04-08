package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	httpapi "pulselounge/internal/http"
	"pulselounge/internal/messages"
)

//go:embed web/dist/*
var embeddedUI embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = buildDatabaseURLFromEnv()
	}

	db, err := initDB(databaseURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database: %v", closeErr)
		}
	}()

	var uiHandler http.Handler
	uiDevServer := os.Getenv("UI_DEV_SERVER")
	if uiDevServer != "" {
		uiHandler, err = spaProxyHandler(uiDevServer)
		if err != nil {
			log.Fatalf("failed to configure UI dev proxy: %v", err)
		}
		log.Printf("proxying UI requests to %s", uiDevServer)
	} else {
		uiFS, fsErr := fs.Sub(embeddedUI, "web/dist")
		if fsErr != nil {
			log.Fatalf("failed to load embedded UI: %v", fsErr)
		}
		uiHandler = spaHandler(uiFS)
	}

	messageRepo := messages.NewPostgresRepository(db)
	messageService := messages.NewService(messageRepo)
	mux := httpapi.NewRouter(uiHandler, messageService)

	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
