package newrelic

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

const (
	DefaultSetupTimeout    = time.Second * 5
	DefaultShutdownTimeout = time.Second * 3
)

var excludeTraceAttributes = []string{
	"appId",
	"duration",
	"entity.guid",
	"entityGuid",
	"id",
	"trace.id",
	"guid",
	"priority",
	"realAgentId",
	"tags.account",
	"tags.accountId",
	"tags.trustedAccountId",
}

type Config struct {
	AppName   string
	License   string
	Debug     bool
	LogOutput io.Writer
	Disabled  bool
}

// TODO: test
func (c Config) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("missing app name")
	}
	if c.License == "" {
		return fmt.Errorf("missing license")
	}

	return nil
}

type Provider struct {
	app *newrelic.Application
	*logger
	*metricRecorder
	*tracer
}

// TODO: remove
func (p *Provider) App() *newrelic.Application {
	return p.app
}

func (p *Provider) Close(ctx context.Context) error {
	timeout := DefaultShutdownTimeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Now().Sub(deadline)
	}

	p.app.Shutdown(timeout)
	return nil
}

func SetupProvider(c Config) (*Provider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	opts := []newrelic.ConfigOption{
		newrelic.ConfigAppName(c.AppName),
		newrelic.ConfigLicense(c.License),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigEnabled(!c.Disabled),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(config *newrelic.Config) {
			// TODO: This doesn't seem to work
			config.TransactionEvents.Attributes.Enabled = true
			config.TransactionEvents.Attributes.Exclude = excludeTraceAttributes

			config.TransactionTracer.Attributes.Enabled = true
			config.TransactionTracer.Attributes.Exclude = excludeTraceAttributes

			config.TransactionTracer.Segments.Attributes.Enabled = true
			config.TransactionTracer.Segments.Attributes.Exclude = excludeTraceAttributes
		},
		// TODO set custom logger
	}

	app, err := newrelic.NewApplication(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create newrelic application")
	}

	if err := app.WaitForConnection(DefaultSetupTimeout); err != nil {
		return nil, errors.Wrap(err, "failed to connect to newrelic")
	}

	if c.LogOutput == nil {
		c.LogOutput = os.Stdout // stderr?
	}

	p := &Provider{
		app:            app,
		logger:         newLogger(app, c.LogOutput, c.Debug),
		metricRecorder: newMetricRecorder(app),
		tracer:         newTracer(app),
	}

	return p, nil
}
