package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort string `envconfig:"GO_UPTIME_MONITOR_HTTP_PORT" default:"8080"`
	GinMode  string `envconfig:"GIN_MODE" default:"release"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `envconfig:"GO_UPTIME_MONITOR_DATABASE_HOST" default:"shared-postgres"`
	Port     string `envconfig:"GO_UPTIME_MONITOR_DATABASE_PORT" default:"5432"`
	User     string `envconfig:"GO_UPTIME_MONITOR_DATABASE_USER" default:"uptimemonitor"`
	DBName   string `envconfig:"GO_UPTIME_MONITOR_DATABASE_NAME" default:"uptimemonitor"`
	Password string `envconfig:"GO_UPTIME_MONITOR_DATABASE_PASSWORD" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return &cfg, nil
}
