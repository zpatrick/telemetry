package newrelic

import (
	"context"
	"io"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/zpatrick/telemetry/tag"
	"golang.org/x/exp/slog"
)

type logger struct {
	writer        logWriter.LogWriter
	defaultLogger *slog.Logger
	debug         bool
	tags          []tag.Tag
}

func newLogger(app *newrelic.Application, w io.Writer, debug bool) *logger {
	writer := logWriter.New(w, app)

	return &logger{
		writer:        writer,
		defaultLogger: newJSONLogger(writer, debug),
		debug:         debug,
	}
}

func (l *logger) log(ctx context.Context, level slog.Level, msg string, tags ...tag.Tag) {
	logger := l.defaultLogger
	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = newJSONLogger(l.writer.WithTransaction(txn), l.debug)
	}

	tags = append(l.tags, append(tag.TagsFromContext(ctx), tags...)...)
	args := make([]any, 0, len(tags)*2)
	tag.Write(tags, func(key string, val any) {
		args = append(args, key, val)
	})

	switch level {
	case slog.LevelDebug:
		logger.DebugCtx(ctx, msg, args...)
	case slog.LevelInfo:
		logger.InfoCtx(ctx, msg, args...)
	case slog.LevelError:
		logger.ErrorCtx(ctx, msg, args...)
	default:
		panic("invalid log level: " + level.String())
	}
}

func (l *logger) Debug(ctx context.Context, msg string, tags ...tag.Tag) {
	l.log(ctx, slog.LevelDebug, msg, tags...)
}

func (l *logger) Info(ctx context.Context, msg string, tags ...tag.Tag) {
	l.log(ctx, slog.LevelInfo, msg, tags...)
}

func (l *logger) Error(ctx context.Context, msg string, tags ...tag.Tag) {
	l.log(ctx, slog.LevelError, msg, tags...)
}

func newJSONLogger(w io.Writer, debug bool) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if debug {
		opts.Level = slog.LevelDebug
	}

	return slog.New(slog.NewJSONHandler(w, opts))
}
