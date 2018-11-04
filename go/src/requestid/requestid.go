package requestid

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type key int

// RequestIDKey specifies the key for the request ID in a request context
const RequestIDKey key = 0

// RequestIDLabel specifies the label for the request ID in a request/response header
const RequestIDLabel = "X-Request-ID"

// Populates a context with a unique request ID if it does not already have one
func newContextWithID(req *http.Request) context.Context {
	ctx := req.Context()
	if FromContext(ctx) != "" {
		return ctx
	}

	return propagateOrGenerate(ctx, req.Header.Get(RequestIDLabel))
}

func propagateOrGenerate(ctx context.Context, id string) context.Context {
	if id == "" {
		id = uuid.NewV1().String()
	}
	return context.WithValue(ctx, RequestIDKey, id)
}

// FromContext grabs the unique ID from the given Context, or returns
// an empty string if it is not found
func FromContext(ctx context.Context) string {
	id, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

// PropagateOrGenerate propagates or generates a unique ID from the request to the response
func PropagateOrGenerate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := newContextWithID(r)
		w.Header().Set(RequestIDLabel, FromContext(ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Generate generates a new unique ID for the context even if a request
// ID previously existed.
func Generate(ctx context.Context) context.Context {
	return propagateOrGenerate(ctx, "")
}
