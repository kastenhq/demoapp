// Package tracing provides functions for configuring tracing and http tracing
// via middleware.
package tracing

import (
	"io"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Middleware returns an http.Handler that uses OpenTracing to trace
func Middleware(next http.Handler) http.Handler {
	// requests that go through it.
	return nethttp.Middleware(opentracing.GlobalTracer(),
		next,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.String()
		}))
}

// Configure configures the Jaeger OpenTracing library so that subsequent
// tracing code operates. If configure is not called, then all other tracing
// code turns into noops.
func Configure(serviceName string) (io.Closer, error) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		return nil, err
	}

	cfg.Reporter.LogSpans = true
	// TODO(ashlie): Change these values to something more sane (ex.
	// probabilistic sampler).
	cfg.Reporter.BufferFlushInterval = 1 * time.Minute
	cfg.Sampler = &jaegercfg.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
	}
	cfg.ServiceName = string(serviceName)
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
