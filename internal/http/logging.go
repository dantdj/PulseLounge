package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

type contextKey string

const (
	loggerContextKey    contextKey = "logger"
	requestIDContextKey contextKey = "request_id"
)

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerContextKey).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}
	return logger
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey).(string)
	return requestID
}

func withRequestLogging(logger *slog.Logger, next http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		requestLogger := logger.With("request_id", requestID)
		ctx := context.WithValue(r.Context(), loggerContextKey, requestLogger)
		ctx = context.WithValue(ctx, requestIDContextKey, requestID)
		r = r.WithContext(ctx)

		w.Header().Set(requestIDHeader, requestID)

		recorder := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		startedAt := time.Now()
		next.ServeHTTP(recorder, r)

		level := slog.LevelInfo
		if recorder.status >= http.StatusInternalServerError {
			level = slog.LevelError
		} else if recorder.status >= http.StatusBadRequest {
			level = slog.LevelWarn
		}

		requestLogger.LogAttrs(
			r.Context(),
			level,
			"http request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", recorder.status),
			slog.Int("response_bytes", recorder.bytesWritten),
			slog.Duration("duration", time.Since(startedAt)),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}

func withPanicRecovery(logger *slog.Logger, next http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				LoggerFromContext(r.Context()).ErrorContext(
					r.Context(),
					"panic while handling request",
					"panic", recovered,
					"method", r.Method,
					"path", r.URL.Path,
				)
				writeJSONError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	written, err := r.ResponseWriter.Write(data)
	r.bytesWritten += written
	return written, err
}
