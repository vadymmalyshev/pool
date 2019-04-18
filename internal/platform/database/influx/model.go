package influx

import "errors"

// DB represent InfluxDB settings for hconsumer uses
type DB struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"password"`
	Name string `yaml:"name"`
}

// Validate returns error if config is not valid
func (db *DB) Validate() error {
	if db.Host == "" {
		return errors.New("influx host must be set")
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
