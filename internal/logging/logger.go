package logging

import (
	"io"
	"log/slog"
	"strings"
)

type Config struct {
	Environment string
	Format      string
	Level       string
	Output      io.Writer
}

func New(config Config) *slog.Logger {
	level := slog.LevelInfo
	switch strings.ToLower(config.Level) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	options := &slog.HandlerOptions{Level: level}
	output := config.Output
	if output == nil {
		output = defaultOutput
	}

	if logFormat(config) == "json" {
		return slog.New(slog.NewJSONHandler(output, options))
	}

	return slog.New(slog.NewTextHandler(output, options))
}

func logFormat(config Config) string {
	format := strings.ToLower(config.Format)
	if format == "json" || format == "text" {
		return format
	}

	if strings.EqualFold(config.Environment, "production") {
		return "json"
	}

	return "text"
}
