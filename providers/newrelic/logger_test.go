package newrelic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/zpatrick/telemetry/tag"
	"github.com/zpatrick/testx/assert"
	"golang.org/x/exp/slog"
)

func TestLoggerLevels(t *testing.T) {
	t.Parallel()

	levels := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelError,
	}

	for _, lvl := range levels {
		t.Run(lvl.String(), func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer(nil)
			logger := newLogger(nil, buf, true)

			var fn func(ctx context.Context, msg string, tags ...tag.Tag)
			switch lvl {
			case slog.LevelDebug:
				fn = logger.Debug
			case slog.LevelInfo:
				fn = logger.Info
			case slog.LevelError:
				fn = logger.Error
			default:
				t.Fatalf("unknown level: %v", lvl)
			}

			ctx := tag.ContextWithTags(context.Background(), tag.New("host", "localhost"))
			fn(ctx, "hello world", tag.New("port", 8080))

			var out struct {
				Message string    `json:"msg"`
				Level   string    `json:"level"`
				Time    time.Time `json:"time"`
				Host    string    `json:"host"`
				Port    int       `json:"port"`
			}
			if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, out.Time.IsZero(), false)
			assert.Equal(t, out.Level, lvl.String())
			assert.Equal(t, out.Message, "hello world")
			assert.Equal(t, out.Host, "localhost")
			assert.Equal(t, out.Port, 8080)
		})
	}
}

func TestLoggerDebugDisabled(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := newLogger(nil, buf, false)

	logger.Debug(context.Background(), "hello world")
	assert.Equal(t, buf.Len(), 0)
}

func TestLoggerWith(t *testing.T) {
	t.Skip("TODO")
}

func TestLoggerE2E(t *testing.T) {
	ensureE2ETestsEnabled(t)

	// TODO: standardize "environemnt" tag.
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

	p.Debug(ctx, "hello debug", tag.New("foo", "bar"))
	p.Info(ctx, "hello info", tag.New("foo", "bar"))
	p.Error(ctx, "hello error", tag.New("foo", "bar"))
}
