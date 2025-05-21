package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host            string        `yaml:"host"`
		Port            string        `yaml:"port"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	} `yaml:"server"`
	BoardService struct {
		Address string `yaml:"address"`
	} `yaml:"board_service"`
	PaymentService struct {
		Address string `yaml:"address"`
	} `yaml:"payment_service"`
	CalendarService struct {
		Address string `yaml:"address"`
	} `yaml:"calendar_service"`
	AuthService struct {
		Address string `yaml:"address"`
	} `yaml:"auth_service"`
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

	return config, nil
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c *Config) GetShutdownTimeoutDuration() time.Duration {
	return c.Server.ShutdownTimeout
}

func (c *Config) GetShutdownTimeoutSeconds() int {
	return int(c.Server.ShutdownTimeout.Seconds())
}

func (c *Config) GetBoardServiceAddr() string {
	return c.BoardService.Address
}

func (c *Config) GetProtectedRoutes() []string {
	return c.Gateway.ProtectedRoutes
}

func (c *Config) GetPaymentServiceAddr() string {
	return c.PaymentService.Address
}

func (c *Config) GetCalendarServiceAddr() string {
	return c.CalendarService.Address
}

func (c *Config) GetAuthServiceAddr() string {
	return c.AuthService.Address
}
