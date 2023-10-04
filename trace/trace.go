package trace

import (
	"context"
	"net/http"

	"github.com/zpatrick/telemetry/tag"
)

type Tracer interface {
	Start(ctx context.Context, name string, tags ...tag.Tag) context.Context
	End(ctx context.Context)

	AddTags(ctx context.Context, tags ...tag.Tag)
	Err(ctx context.Context, err error)

	TraceHandler(h http.Handler) http.Handler
	RoundTripper(r http.RoundTripper) http.RoundTripper
}
