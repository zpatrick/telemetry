package newrelic

import (
	"context"
	"net/http"
	"reflect"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/zpatrick/telemetry/tag"
)

const (
	outgoingHTTPRequestTransactionName = "outgoing_http_request"
	incomingHTTPRequestTransactionName = "incoming_http_request"
)

type tracer struct {
	app *newrelic.Application
}

func newTracer(app *newrelic.Application) *tracer {
	return &tracer{
		app: app,
	}
}

func (t *tracer) Start(ctx context.Context, name string, tags ...tag.Tag) context.Context {
	if txn := newrelic.FromContext(ctx); txn != nil {
		return t.startSegment(ctx, txn, name, tags...)
	}

	return t.startTransaction(ctx, name, tags...)
}

func (t *tracer) startSegment(ctx context.Context, txn *newrelic.Transaction, name string, tags ...tag.Tag) context.Context {
	seg := txn.StartSegment(name)
	t.addAttributes(ctx, seg, tags...)

	return contextWithSegment(ctx, seg)
}

func (t *tracer) startTransaction(ctx context.Context, name string, tags ...tag.Tag) context.Context {
	txn := t.app.StartTransaction(name)
	t.addAttributes(ctx, txn, tags...)

	return newrelic.NewContext(ctx, txn)
}

func (t *tracer) AddTags(ctx context.Context, tags ...tag.Tag) {
	if seg := segmentFromContext(ctx); seg != nil {
		t.addAttributes(ctx, seg, tags...)
		return
	}

	if txn := newrelic.FromContext(ctx); txn != nil {
		t.addAttributes(ctx, txn, tags...)
		return
	}
}

func (t *tracer) addAttributes(ctx context.Context, a interface{ AddAttribute(string, interface{}) }, tags ...tag.Tag) {
	tags = append(tag.TagsFromContext(ctx), tags...)
	tag.Write(tags, func(key string, value interface{}) {
		a.AddAttribute(key, value)
	})
}

func (t *tracer) End(ctx context.Context) {
	if seg := segmentFromContext(ctx); seg != nil {
		seg.End()
		return
	}

	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.End()
		return
	}
}

func (t *tracer) Err(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return err
	}

	// TODO: first-class error types, categories, tags, etc.
	var tags map[string]interface{}
	if t, ok := err.(interface{ Tags() []tag.Tag }); ok {
		tag.Write(t.Tags(), func(key string, value interface{}) {
			tags[key] = value
		})
	}

	txn.NoticeError(newrelic.Error{
		Message:    err.Error(),
		Class:      reflect.TypeOf(err).String(),
		Attributes: tags,
	})

	return err
}

func (t *tracer) TraceHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := t.app.StartTransaction(incomingHTTPRequestTransactionName)
		t.addAttributes(r.Context(), txn)
		txn.SetWebRequestHTTP(r)
		defer txn.End()

		h.ServeHTTP(txn.SetWebResponse(w), newrelic.RequestWithTransactionContext(r, txn))
	})
}

func (t *tracer) RoundTripper(rt http.RoundTripper) http.RoundTripper {
	wrapped := newrelic.NewRoundTripper(rt)
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if txn := newrelic.FromContext(req.Context()); txn != nil {
			req = newrelic.RequestWithTransactionContext(req, txn)
			return wrapped.RoundTrip(req)
		}

		txn := t.app.StartTransaction(outgoingHTTPRequestTransactionName)
		t.addAttributes(req.Context(), txn)
		defer txn.End()

		req = newrelic.RequestWithTransactionContext(req, txn)
		return wrapped.RoundTrip(req)
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type segmentKey int

func contextWithSegment(ctx context.Context, seg *newrelic.Segment) context.Context {
	return context.WithValue(ctx, segmentKey(0), seg)
}

func segmentFromContext(ctx context.Context) *newrelic.Segment {
	seg, _ := ctx.Value(segmentKey(0)).(*newrelic.Segment)
	return seg
}
