package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/matt-FFFFFF/otelcli/internal/appinsights"
)

func main() {
	ctx := context.Background()
	conStr := flag.String("c", "", "The connection string for the Azure Monitor exporter")
	flag.Parse()

	if *conStr == "" {
		fmt.Println("Connection string is required")
		os.Exit(1)
	}

	client, err := appinsights.NewClientFromConnectionString(ctx, *conStr)
	if err != nil {
		fmt.Printf("Error creating Application Insights client: %v\n", err)
		os.Exit(1)
	}
	defer client.Shutdown()
	client.Log("Hello, world!", nil)
}
