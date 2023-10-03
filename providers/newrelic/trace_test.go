package newrelic

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zpatrick/telemetry/tag"
	"github.com/zpatrick/testx/assert"
)

// Test Scenarios:

// Receiving a trace from a client.
// Sending a trace to a client.

func TestTracerE2E_basic(t *testing.T) {
	ctx, p := newTestProvider(t)
	ctx = p.Start(ctx, "TracerTestBasic", tag.New("alpha", 1))
	defer p.End(ctx)

	p.AddTags(ctx, tag.New("beta", 2))
}

func TestTracerE2E_error(t *testing.T) {
	ctx, p := newTestProvider(t)
	ctx = p.Start(ctx, "TracerTestError")
	defer p.End(ctx)

	// TODO: test attributes, class, etc.
	err := errors.New("test error")
	p.Err(ctx, err)
}

func TestTracerE2E_nested(t *testing.T) {
	ctx, p := newTestProvider(t)
	ctx = p.Start(ctx, "TracerTestNested: root", tag.New("alpha", 1))
	defer p.End(ctx)

	inner := func(ctx context.Context) {
		ctx = p.Start(ctx, "TracerTestNested: leaf", tag.New("beta", 2))
		defer p.End(ctx)
	}

	inner(ctx)
}

func TestTracerE2E_http(t *testing.T) {
	ctx, p := newTestProvider(t)
	ctx = p.Start(ctx, "TracerTestHTTP: root")
	defer p.End(ctx)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := p.Start(r.Context(), "TracerTestHTTP: handler")
		defer p.End(ctx)

		w.Write([]byte("hello world"))
	})

	svr := httptest.NewServer(p.TraceHandler(handler))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svr.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	client.Transport = p.RoundTripper(client.Transport)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
