package trace

import (
	"context"

	"github.com/zpatrick/telemetry/tag"
)

type Tracer interface {
	Start(ctx context.Context, name string, tags ...tag.Tag) context.Context
	AddTags(ctx context.Context, tags ...tag.Tag) // context.Context?
	Finish(ctx context.Context)
}

func Error(ctx context.Context, t Tracer, err error) error {
	// Adds error metadata to existing trace.
	return err
}
