package mysql

import "errors"

// DB represents MySQL db config
type DB struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	SSLMode bool   `yaml:"sslmode"`
	User    string `yaml:"user"`
	Pass    string `yaml:"password"`
	Name    string `yaml:"name"`
	Log     bool   `yaml:"log"`
}

// Validate checks that Mysql config is valid
func (db *DB) Validate() error {
	if db.Host == "" {
		return errors.New("mysql host must be set")
	}

	if db.Port < 1025 || db.Port > 65535 {
		return errors.New("influx port must be set")
	}

	if db.User == "" {
		return errors.New("influx user must be set")
	}

	if db.Pass == "" {
		return errors.New("influx user must be set")
	}

	if db.Name == "" {
		return errors.New("influx db name must be set")
	}

	return nil
}
