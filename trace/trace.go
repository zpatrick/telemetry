package trace

import (
	"context"
	"net/http"

	"github.com/zpatrick/telemetry/tag"
)

type Tracer interface {
	Start(ctx context.Context, name string, tags ...tag.Tag) context.Context
	End(ctx context.Context)

	AddTags(ctx context.Context, tags ...tag.Tag) // context.Context?
	Err(ctx context.Context, err error)

	TraceHandler(h http.Handler) http.Handler
	RoundTripper(r http.RoundTripper) http.RoundTripper
}

func Error(ctx context.Context, t Tracer, err error) error {
	// Adds error metadata to existing trace.
	return err
}

// TODO: add a middleware that adds trace metadata to the context.
// TODO: add a http client round tripper which adds middleware.
