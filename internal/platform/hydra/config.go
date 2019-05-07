package hydra

import "errors"

// Config represents hydra settings
type Config struct {
	API    string `yaml:"api"`
	Admin  string `yaml:"admin"`
	Docker string `yaml:"docker"`
}

// Validate checks hydra settings
func (c Config) Validate() error {
	if c.API == "" {
		return errors.New("hydra public url is required")
	}

	if c.Admin == "" {
		return errors.New("hydra admin url is required")
	}

	return nil
}
