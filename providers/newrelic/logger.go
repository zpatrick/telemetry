package newrelic

import (
	"context"
	"io"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/zpatrick/telemetry/log"
	"github.com/zpatrick/telemetry/tag"
	"golang.org/x/exp/slog"
)

type logger struct {
	logger *slog.Logger
	tags   []tag.Tag
}

func newLogger(app *newrelic.Application, w io.Writer, debug bool) *logger {
	writer := logWriter.New(w, app)
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if debug {
		opts.Level = slog.LevelDebug
	}

	return &logger{
		logger: slog.New(slog.NewJSONHandler(&writer, opts)),
	}
}

func (l *logger) log(ctx context.Context, fn func(ctx context.Context, msg string, args ...any), msg string, tags ...tag.Tag) {
	tags = append(l.tags, append(tag.TagsFromContext(ctx), tags...)...)

	args := make([]any, 0, len(tags)*2)
	tag.Write(tags, func(key string, val any) {
		args = append(args, key, val)
	})

	fn(ctx, msg, args...)
}

func (l *logger) Debug(ctx context.Context, msg string, tags ...tag.Tag) {
	l.log(ctx, l.logger.DebugContext, msg, tags...)
}

func (l *logger) Info(ctx context.Context, msg string, tags ...tag.Tag) {
	l.log(ctx, l.logger.InfoContext, msg, tags...)
}

func (l *logger) Error(ctx context.Context, msg string, tags ...tag.Tag) {
	l.log(ctx, l.logger.ErrorContext, msg, tags...)
}

// TODO: newrelic has primatives to add transactional logging. Should we tap into that?
// We could have a method in the telemetry pacakge which tries to write them up?
// func TraceLogger(ctx context.Context, p Provider) log.Logger {
// 	if p, ok := p.(TraceLogger); ok {
// 		return p.TraceLogger(ctx)
// 	}

// 	return p.Logger
// }

func (l *logger) With(tags ...tag.Tag) log.Logger {
	// TODO: l.logger.With(tags...)
	return &logger{logger: l.logger, tags: append(l.tags, tags...)}
}
