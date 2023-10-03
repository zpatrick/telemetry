package newrelic

import (
	"math/rand"
	"testing"
	"time"
)

func TestMetricE2E(t *testing.T) {
	ctx, p := newTestProvider(t)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		p.Count(ctx, "testEvent", 1)
		p.Guage(ctx, "test_request_duration_ms", rand.Int63n(2000))
	}
}
