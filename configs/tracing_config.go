package configs

import (
	"go.uber.org/zap"
)

const (
	newRelicProvider = "new-relic"
)

type tracingConfig struct {
	// Enabled specifies that the tracing agent will be created and used.
	Enabled bool

	// Provider is the supported providers. Only NewRelic is supported for now.
	Provider string

	// NewRelic provider options.
	NewRelic newRelicConfig
}

type newRelicConfig struct {
	// AppName is the name of the application on NewRelic.
	AppName string

	// LicenseKey is a 40 character length secret.
	LicenseKey string

	// DistributedTracerEnabled should be enabled when DistributedTracing feature wanted to use.
	DistributedTracerEnabled bool
}

func (t *tracingConfig) validateProvider() {
	if !t.Enabled {
		return
	}

	switch t.Provider {
	case newRelicProvider:
	default:
		t.Enabled = false
		zap.L().Warn("Tracing is set as disabled due to invalid provider name.", zap.String("providerName", t.Provider))
	}
}
