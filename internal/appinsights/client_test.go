package appinsights

import (
	"context"
	"testing"
)

func TestNewLogs(t *testing.T) {
	_, err := NewClientFromConnectionString(context.Background(), "InstrumentationKey=00000000-0000-0000-0000-000000000000;IngestionEndpoint=https://region-0.in.applicationinsights.azure.com/;LiveEndpoint=https://region.livediagnostics.monitor.azure.com/;ApplicationId=00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
}
