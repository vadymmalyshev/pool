package database

import (
	"errors"
	"fmt"
)

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

	if c.Name == "" {
		return errors.New("database name is required")
	}

	return nil
}

// Connection returns db connection string
func (c Config) Connection() string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "enable"
	}
	return fmt.Sprintf("host=%s port=%d sslmode=%s user=%s dbname=%s password=%s ",
		c.Host,
		c.Port,
		ssl,
		c.User,
		c.Name,
		c.Pass)
}

// Connection returns db connection string
func (c Config) MySQLConnection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name)
}