package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/matt-FFFFFF/otelcli/internal/appinsights"
)

func main() {
	ctx := context.Background()
	ingestUrl := flag.String("u", "", "The ingestion URL, e.g. https://region-0.in.applicationinsights.azure.com/v2.1/track")
	instrumentationKey := flag.String("i", "", "The instrumentation key, e.g. 00000000-0000-0000-0000-000000000000")
	numEvents := flag.Int("n", 1, "The number of events to log")
	flag.Parse()

	if *ingestUrl == "" || *instrumentationKey == "" {
		fmt.Println("Ingestion url and instrumentation key are required")
		os.Exit(1)
	}

	client, err := appinsights.NewClient(ctx, *instrumentationKey, *ingestUrl)
	if err != nil {
		fmt.Printf("Error creating Application Insights client: %v\n", err)
		os.Exit(1)
	}
	defer client.Shutdown()

	// Log messages
	for i := range *numEvents {
		fmt.Printf("Logging event... %d\n", i+1)
		client.Log(eventName(), eventProperties())
	}
	fmt.Printf("Quitting...\n")
}

func eventProperties() map[string]string {
	return map[string]string{
		"module_source":   "registry.terraform.io/Azure/avm-utl-regions/azurerm",
		"module_version":  "1.0.0",
		"subscription_id": "00000000-0000-0000-0000-000000000000",
		"tenant_id":       "00000000-0000-0000-0000-000000000000",
		"resource_id":     "00000000-0000-0000-0000-000000000000",
	}
}

func eventName() string {
	events := []string{
		"read",
		"create",
		"delete",
	}
	return events[rand.Intn(len(events))]
}
