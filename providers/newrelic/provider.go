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
	}

	return p, nil
}
