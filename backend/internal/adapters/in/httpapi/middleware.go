package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

const requestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode != 0 {
		return
	}
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusResponseWriter) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := newRequestID()
		w.Header().Set(requestIDHeader, requestID)
		contextWithRequestID := context.WithValue(r.Context(), requestIDContextKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(contextWithRequestID))
	})
}

func withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		responseWriter := &statusResponseWriter{ResponseWriter: w}
		next.ServeHTTP(responseWriter, r)

		statusCode := responseWriter.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		slog.Info("http request completed",
			"request_id", requestIDFromContext(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", statusCode,
			"duration_ms", time.Since(startedAt).Milliseconds(),
		)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recoveredValue := recover(); recoveredValue != nil {
				slog.Error("panic recovered while processing request",
					"request_id", requestIDFromContext(r.Context()),
					"method", r.Method,
					"path", r.URL.Path,
					"panic", fmt.Sprint(recoveredValue),
					"stack_trace", string(debug.Stack()),
				)

				if responseWriter, ok := w.(*statusResponseWriter); !ok || responseWriter.statusCode == 0 {
					writeErrorResponse(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func logRequestRejection(r *http.Request, errorCode string, operation string, operandCount int) {
	attributes := []any{
		"request_id", requestIDFromContext(r.Context()),
		"error_code", errorCode,
		"operand_count", operandCount,
	}
	if operation != "" {
		attributes = append(attributes, "operation", operation)
	}
	slog.Warn("calculator request rejected", attributes...)
}

func requestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}

func newRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err == nil {
		return hex.EncodeToString(bytes)
	}
	return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
}
