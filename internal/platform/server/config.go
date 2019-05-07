package server

import (
	"errors"
	"fmt"
)

// Config holds information necessary for running a server
type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Validate checks that the configuration is valid
func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("host is required")
	}

	if c.Port == 0 {
		return errors.New("port is required")
	}

	if c.Port < 1025 || c.Port > 65535 {
		return errors.New("port is wrong")
	}

	return nil
}

// Addr returns server address
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
