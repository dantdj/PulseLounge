package main

import (
	"log/slog"
	"net/http"
	"os"

	uiassets "pulselounge"
	"pulselounge/internal/channels"
	httpapi "pulselounge/internal/http"
	"pulselounge/internal/logging"
	"pulselounge/internal/media"
	"pulselounge/internal/messages"
)

func main() {
	logger := logging.New(logging.Config{
		Environment: appEnvironment(),
		Format:      os.Getenv("LOG_FORMAT"),
		Level:       os.Getenv("LOG_LEVEL"),
	})
	slog.SetDefault(logger)

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
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("failed to close database", "error", closeErr)
		}
	}()

	connectionString := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	if connectionString == "" {
		logger.Error("AZURE_STORAGE_CONNECTION_STRING environment variable is required")
		os.Exit(1)
	}
	publicBaseURL := os.Getenv("MEDIA_PUBLIC_BASE_URL")
	if publicBaseURL == "" {
		logger.Error("MEDIA_PUBLIC_BASE_URL environment variable is required")
		os.Exit(1)
	}
	containerName := os.Getenv("MEDIA_CONTAINER_NAME")
	if containerName == "" {
		logger.Error("MEDIA_CONTAINER_NAME environment variable is required")
		os.Exit(1)
	}

	var uiHandler http.Handler
	uiDevServer := os.Getenv("UI_DEV_SERVER")
	if uiDevServer != "" {
		uiHandler, err = spaProxyHandler(uiDevServer)
		if err != nil {
			logger.Error("failed to configure UI dev proxy", "error", err)
			os.Exit(1)
		}
		logger.Info("proxying UI requests", "target", uiDevServer)
	} else {
		uiFS, fsErr := uiassets.FS()
		if fsErr != nil {
			logger.Error("failed to load embedded UI", "error", fsErr)
			os.Exit(1)
		}
		uiHandler = spaHandler(uiFS)
	}

	messageRepo := messages.NewPostgresRepository(db)
	messageService := messages.NewService(messageRepo)
	channelRepo := channels.NewPostgresRepository(db)
	channelService := channels.NewService(channelRepo)
	blobStore, err := media.NewAzureBlobStore(connectionString, containerName, publicBaseURL)
	if err != nil {
		logger.Error("failed to configure blob storage", "error", err)
		os.Exit(1)
	}
	mux := httpapi.NewRouterWithLogger(uiHandler, messageService, channelService, blobStore, logger)

	addr := ":" + port
	logger.Info("server listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func appEnvironment() string {
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}
	return os.Getenv("GO_ENV")
}
