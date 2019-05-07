package database

import (
	"errors"
)

// DBConfiger is the main interface to db providers
type DBConfiger interface {
	Validate() error
	Connect() (interface{}, error)
}

// Config represents db config struct
type Config struct {
	Host      string
	Port      int
	EnableSSL bool
	User      string
	Pass      string
	Name      string
	EnableLog bool
}

// Validate returns error if config is not valid
func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required")
	}

	if c.Port == 0 {
		return errors.New("database port is required")
	}

	if c.Port < 1025 || c.Port > 65535 {
		return errors.New("database port is invalid")
	}

	if c.Name == "" {
		return errors.New("database name is required")
	}

	return nil
}
