package newrelic

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/zpatrick/telemetry/tag"
)

type metricRecorder struct {
	app *newrelic.Application
}

func newMetricRecorder(app *newrelic.Application) *metricRecorder {
	return &metricRecorder{
		app: app,
	}
}

func (r *metricRecorder) record(ctx context.Context, name string, value int64) {
	params := map[string]interface{}{"value": value}
	tag.Write(tag.TagsFromContext(ctx), func(key string, val any) {
		params[key] = val
	})

	r.app.RecordCustomEvent(name, params)
}

func (r *metricRecorder) Count(ctx context.Context, name string, value int64) {
	r.record(ctx, name, value)
}

func (r *metricRecorder) Guage(ctx context.Context, name string, value int64) {
	r.record(ctx, name, value)
}
