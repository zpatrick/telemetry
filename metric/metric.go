package metric

import (
	"context"

	"github.com/zpatrick/telemetry/tag"
)

type Recorder interface {
	Count(ctx context.Context, name string, value int64, tags ...tag.Tag)
	Gauge(ctx context.Context, name string, value float64, tags ...tag.Tag)
}
