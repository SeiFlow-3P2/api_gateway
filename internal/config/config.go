package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Addr            string `yaml:"addr"`
		ShutdownTimeout int    `yaml:"shutdown_timeout"`
	} `yaml:"server"`
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

	return config, nil
}

func (c *Config) GetServerAddr() string {
	return c.Server.Addr
}

func (c *Config) GetShutdownTimeout() int {
	return c.Server.ShutdownTimeout
}
