package main

import (
	"log"
	"net/http"
	"os"

	uiassets "pulselounge"
	"pulselounge/internal/channels"
	httpapi "pulselounge/internal/http"
	"pulselounge/internal/media"
	"pulselounge/internal/messages"
)

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

	connectionString := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("AZURE_STORAGE_CONNECTION_STRING environment variable is required")
	}
	publicBaseURL := os.Getenv("MEDIA_PUBLIC_BASE_URL")
	if publicBaseURL == "" {
		log.Fatal("MEDIA_PUBLIC_BASE_URL environment variable is required")
	}
	containerName := os.Getenv("MEDIA_CONTAINER_NAME")
	if containerName == "" {
		log.Fatal("MEDIA_CONTAINER_NAME environment variable is required")
	}

	var uiHandler http.Handler
	uiDevServer := os.Getenv("UI_DEV_SERVER")
	if uiDevServer != "" {
		uiHandler, err = spaProxyHandler(uiDevServer)
		if err != nil {
			log.Fatalf("failed to configure UI dev proxy: %v", err)
		}
		log.Printf("proxying UI requests to %s", uiDevServer)
	} else {
		uiFS, fsErr := uiassets.FS()
		if fsErr != nil {
			log.Fatalf("failed to load embedded UI: %v", fsErr)
		}
		uiHandler = spaHandler(uiFS)
	}

	messageRepo := messages.NewPostgresRepository(db)
	messageService := messages.NewService(messageRepo)
	channelRepo := channels.NewPostgresRepository(db)
	channelService := channels.NewService(channelRepo)
	blobStore, err := media.NewAzureBlobStore(connectionString, containerName, publicBaseURL)
	if err != nil {
		log.Fatalf("failed to configure blob storage: %v", err)
	}
	mux := httpapi.NewRouter(uiHandler, messageService, channelService, blobStore)

	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
