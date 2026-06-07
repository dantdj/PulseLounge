package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewUsesJSONForProduction(t *testing.T) {
	var logs bytes.Buffer
	logger := New(Config{Environment: "production", Output: &logs})

	logger.Info("started", "component", "server")

	if !strings.HasPrefix(logs.String(), "{") {
		t.Fatalf("expected JSON log, got %q", logs.String())
	}
	if !strings.Contains(logs.String(), `"component":"server"`) {
		t.Fatalf("expected structured JSON field, got %q", logs.String())
	}
}

func TestNewUsesTextOutsideProduction(t *testing.T) {
	var logs bytes.Buffer
	logger := New(Config{Environment: "development", Output: &logs})

	logger.Info("started", "component", "server")

	if strings.HasPrefix(logs.String(), "{") {
		t.Fatalf("expected text log, got %q", logs.String())
	}
	if !strings.Contains(logs.String(), "component=server") {
		t.Fatalf("expected structured text field, got %q", logs.String())
	}
}

func TestNewAllowsFormatOverride(t *testing.T) {
	var logs bytes.Buffer
	logger := New(Config{Environment: "production", Format: "text", Output: &logs})

	logger.Info("started", "component", "server")

	if strings.HasPrefix(logs.String(), "{") {
		t.Fatalf("expected text log from override, got %q", logs.String())
	}
}
