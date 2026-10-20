// Package trace provides trace ID propagation through context, HTTP, and gRPC.
package trace

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type key int

const traceKey key = iota

// NewID generates a UUIDv7 trace ID.
func NewID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}

// NewContext returns a context with the given trace ID.
func NewContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey, traceID)
}

// FromContext extracts the trace ID, or "" if none is set.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(traceKey).(string)
	return id
}

// HTTPMiddleware wraps an http.Handler, injecting a trace ID into every request.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = NewID()
		}
		w.Header().Set("X-Trace-Id", traceID)
		ctx := NewContext(r.Context(), traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}