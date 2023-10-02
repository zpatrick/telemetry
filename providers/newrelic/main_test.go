package newrelic

import (
	"flag"
	"os"
	"testing"
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

func TestMain(m *testing.M) {
	flag.BoolVar(&e2eTestsEnabled, "e2e", false, "enable e2e tests")
	flag.Parse()

	os.Exit(m.Run())
}
