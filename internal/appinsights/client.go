package appinsights

import (
	"context"
	"fmt"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

const (
	// MaxBatchSize is the maximum number of telemetry items that can be sent in one call to the data collector.
	MaxBatchSize = 4096
)

type (
	// Client is an interface for logging telemetry
	Client interface {
		Log(msg string, properties map[string]string)    // Log sends a telemetry message to Azure Application Insights.
		LogNow(msg string, properties map[string]string) // LogNow sends a telemetry message to Azure Application Insights and flushes the channel.
		Shutdown()                                       // Shutdown closes the telemetry channel and waits for all telemetry to be sent.
	}

	// client is a concrete implementation of the Client interface for Azure Application Insights.
	client struct {
		telemetryClient appinsights.TelemetryClient
		shutdownFunc    ShutdownFunc
	}

	// ShutdownFunc is a function type that represents a shutdown function for the logger.
	// It is used to ensure that the telemetry channel is closed properly.
	ShutdownFunc func()
)

func (c *client) Log(msg string, properties map[string]string) {
	event := appinsights.NewEventTelemetry(msg)
	event.Properties = properties
	c.telemetryClient.Track(event)
}

func (c *client) LogNow(msg string, properties map[string]string) {
	c.Log(msg, properties)
	c.telemetryClient.Channel().Flush()
}

func (c *client) Shutdown() {
	c.shutdownFunc()
}

func NewClientFromConnectionString(ctx context.Context, connectionString string) (Client, error) {
	vals, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	return NewClient(ctx, vals.InstrumentationKey, vals.IngestionURL)
}

// NewClient creates a new Application Insights client with the given instrumentation key and ingestion URL.
// The ingestion URL is optional; if not provided, the default ingestion URL will be used.
// The ingestion URL should be in the format "https://something.in.applicationinsights.azure.com/v2.1/track".
func NewClient(ctx context.Context, instrumentationKey, ingestionURL string) (Client, error) {
	if instrumentationKey == "" {
		return nil, fmt.Errorf("instrumentation key cannot be empty")
	}

	telemetryConfig := appinsights.NewTelemetryConfiguration(instrumentationKey)
	if ingestionURL != "" {
		telemetryConfig.EndpointUrl = ingestionURL
	}

	// Configure how many items can be sent in one call to the data collector:
	telemetryConfig.MaxBatchSize = MaxBatchSize

	// Configure the maximum delay before sending queued telemetry:
	telemetryConfig.MaxBatchInterval = 2 * time.Second

	tc := appinsights.NewTelemetryClientFromConfig(telemetryConfig)
	shut := ShutdownFunc(func() {
		cl := tc.Channel().Close()
		<-cl
	})
	return &client{
		telemetryClient: tc,
		shutdownFunc:    shut,
	}, nil
}
