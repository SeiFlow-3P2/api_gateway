package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Name            string        `yaml:"name"`
		Host            string        `yaml:"host"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	} `yaml:"server"`
	Gateway struct {
		ProtectedRoutes []string `yaml:"protected_routes"`
	} `yaml:"gateway"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if config.Server.Name == "" {
		config.Server.Name = "api_gateway"
	}

	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}

	if config.Server.ShutdownTimeout == 0 {
		config.Server.ShutdownTimeout = 10 * time.Second
	}

	if len(config.Gateway.ProtectedRoutes) == 0 {
		return nil, fmt.Errorf("gateway.protected_routes is required")
	}

	return config, nil
}

func (c *Config) GetServerName() string {
	return c.Server.Name
}

func (c *Config) GetShutdownTimeoutDuration() time.Duration {
	return c.Server.ShutdownTimeout
}

func (c *Config) GetShutdownTimeoutSeconds() int {
	return int(c.Server.ShutdownTimeout.Seconds())
}

func (c *Config) GetProtectedRoutes() []string {
	return c.Gateway.ProtectedRoutes
}
