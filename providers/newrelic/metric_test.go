package newrelic

import (
	"context"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/zpatrick/telemetry/tag"
)

func TestMetricE2E(t *testing.T) {
	ensureE2ETestsEnabled(t)

	ctx := tag.ContextWithTags(context.Background(), tag.New("environment", "test"))
	cfg := Config{
		AppName:   "test-app",
		License:   newRelicLicense,
		LogOutput: io.Discard,
		Debug:     true,
	}

	p, err := SetupProvider(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close(ctx)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		p.Count(ctx, "testEvent", 1)
		p.Guage(ctx, "test_request_duration_ms", rand.Int63n(2000))
	}
}
