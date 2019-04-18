package hydraclient

import "errors"

// Config represents hydra client credentials
type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	CallbackURL  string `yaml:"callback"`
}

// Validate checks hydra client settings
func (c Config) Validate() error {
	if c.ClientID == "" {
		return errors.New("client id is required")
	}

	if c.ClientSecret == "" {
		return errors.New("client secret is required")
	}

	if c.CallbackURL == "" {
		return errors.New("callback url is required")
	}

	return nil
}
