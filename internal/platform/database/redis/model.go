package redis

import "errors"

// DB represents Redis config
type DB struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"db"`
	Pass string `yaml:"password"`
}

// Validate checks that Mysql config is valid
func (db *DB) Validate() error {
	if db.Host == "" {
		return errors.New("redis host must be set")
	}
	if db.Port < 1025 || db.Port > 65535 {
		return errors.New("redis port must be set")
	}

	if db.Pass == "" {
		return errors.New("redis user must be set")
	}

	if db.Name == "" {
		return errors.New("redis db name must be set")
	}

	return nil
}
