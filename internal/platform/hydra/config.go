package hydra

import "errors"

type Config struct {
	APIUrl   string
	AdminURL string
}

func (c Config) Validate() error {
	if c.APIUrl == "" {
		return errors.New("hydra public url is required")
	}

	if c.AdminURL == "" {
		return errors.New("hydra admin url is required")
	}

	return nil
}
