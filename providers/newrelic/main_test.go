package newrelic

import (
	"context"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/zpatrick/telemetry/tag"
)

var (
	e2eTestsEnabled bool
	newRelicLicense string = os.Getenv("NEWRELIC_LICENSE")
)

func ensureE2ETestsEnabled(t *testing.T) {
	if !e2eTestsEnabled {
		t.Skip("e2e tests disabled")
	}

	if newRelicLicense == "" {
		t.Fatal("NEWRELIC_LICENSE not set")
	}
}

func newTestProvider(t *testing.T) (context.Context, *Provider) {
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
	t.Cleanup(func() { p.Close(ctx) })

	return ctx, p
}

func TestMain(m *testing.M) {
	flag.BoolVar(&e2eTestsEnabled, "e2e", false, "enable e2e tests")
	flag.Parse()

	os.Exit(m.Run())
}
