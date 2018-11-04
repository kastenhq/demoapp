// Package transport provides functions for getting an HTTP transport capable of
// performing tracing.
package transport

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracingTransport struct {
	transport http.RoundTripper
}

func (t *tracingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	span, ctx2 := opentracing.StartSpanFromContext(ctx, "HTTP Request")
	defer span.Finish()
	span.SetTag(string(ext.SpanKind), ext.SpanKindRPCServerEnum)
	span.SetTag(string(ext.HTTPMethod), r.Method)
	span.SetTag(string(ext.HTTPUrl), r.URL)

	r.WithContext(ctx2)

	resp, err := t.transport.RoundTrip(r)

	span.SetTag(string(ext.HTTPStatusCode), resp.StatusCode)

	if err != nil {
		span.SetTag(string(ext.Error), true)
		return nil, err
	}

	return resp, err
}

func NewTracingTransport(rt http.RoundTripper) (http.RoundTripper) {
	return &tracingTransport{transport: rt}
}
