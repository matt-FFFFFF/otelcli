// Copyright The OpenTelemetry Authors

package appinsights

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

type ConnectionVars struct {
	InstrumentationKey string
	IngestionURL       string
}

const (
	ApplicationInsightsConnectionString = "APPLICATIONINSIGHTS_CONNECTION_STRING"
	DefaultIngestionEndpoint            = "https://dc.services.visualstudio.com/"
	IngestionEndpointKey                = "IngestionEndpoint"
	InstrumentationKey                  = "InstrumentationKey"
	ConnectionStringMaxLength           = 4096
)

func parseConnectionString(connectionString string) (*ConnectionVars, error) {
	// First, try to get the connection string from the environment variable
	if connectionString == "" {
		connectionString = os.Getenv(ApplicationInsightsConnectionString)
	}
	// If the connection string is still empty, return an error
	if connectionString == "" {
		return nil, fmt.Errorf("ConnectionString cannot be empty")
	}

	connectionVars := &ConnectionVars{}

	if len(connectionString) > ConnectionStringMaxLength {
		return nil, fmt.Errorf("ConnectionString exceeds maximum length of %d characters", ConnectionStringMaxLength)
	}

	pairs := strings.Split(connectionString, ";")
	values := make(map[string]string)
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid format for connection string: %s", pair)
		}

		key, value := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		if key == "" {
			return nil, fmt.Errorf("key cannot be empty")
		}
		values[key] = value
	}

	var ok bool
	if connectionVars.InstrumentationKey, ok = values[InstrumentationKey]; !ok || connectionVars.InstrumentationKey == "" {
		return nil, fmt.Errorf("%s is required", InstrumentationKey)
	}

	var ingestionEndpoint string
	if ingestionEndpoint, ok = values[IngestionEndpointKey]; !ok || ingestionEndpoint == "" {
		ingestionEndpoint = DefaultIngestionEndpoint
	}

	connectionVars.IngestionURL = getIngestionURL(ingestionEndpoint)

	return connectionVars, nil
}

func getIngestionURL(ingestionEndpoint string) string {
	ingestionURL, err := url.Parse(ingestionEndpoint)
	if err != nil {
		ingestionURL, _ = url.Parse(DefaultIngestionEndpoint)
	}

	ingestionURL.Path = path.Join(ingestionURL.Path, "/v2.1/track")
	return ingestionURL.String()
}
