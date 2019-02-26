package server

import (
	"errors"
	"fmt"
)

// Config holds information necessary for running a server
type Config struct {
	Host     string
	Port     int
	CertFile string
	KeyFile  string
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

	if c.CertFile == "" {
		return errors.New(".pem cert file is required")
	}

	if c.KeyFile == "" {
		return errors.New(".key cert file is required")
	}

	return nil
}

// Addr returns server address
func (c Config) Addr() string {
	return fmt.Sprintf("https://%s:%d/", c.Host, c.Port)
}
