package hydraclient

import "errors"

type Config struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

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
