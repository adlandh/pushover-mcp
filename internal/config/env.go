package config

import (
	"fmt"
	"time"

	"github.com/adlandh/pushover-mcp/internal/adapters"
	"github.com/caarlos0/env/v11"
)

type EnvConfig struct {
	Pushover adapters.Config
	Timeout  time.Duration
}

type rawEnvConfig struct {
	PushoverAPIToken string        `env:"PUSHOVER_API_TOKEN,notEmpty"`
	PushoverUserKey  string        `env:"PUSHOVER_USER_KEY,notEmpty"`
	PushoverAPIURL   string        `env:"PUSHOVER_API_URL"`
	PushoverTimeout  time.Duration `env:"PUSHOVER_TIMEOUT" envDefault:"15s"`
}

func FromEnv() (EnvConfig, error) {
	var raw rawEnvConfig
	if err := env.Parse(&raw); err != nil {
		return EnvConfig{}, fmt.Errorf("parse env: %w", err)
	}

	return EnvConfig{Pushover: adapters.Config{
		APIToken: raw.PushoverAPIToken,
		UserKey:  raw.PushoverUserKey,
		APIURL:   raw.PushoverAPIURL,
	},
		Timeout: raw.PushoverTimeout,
	}, nil
}
