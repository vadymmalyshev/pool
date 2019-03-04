package database

import (
	"testing"
)

func TestValidConfig(t *testing.T) {

	cases := []struct {
		name         string
		config       Config
		failExpected bool
	}{
		{
			name: "Database valid config",
			config: Config{
				Host: "localhost",
				Port: 54320,
				Name: "dbname",
			},
			failExpected: false,
		},
		{
			name: "Database with no name config",
			config: Config{
				Host: "localhost",
				Port: 54320,
			},
			failExpected: true,
		},
		{
			name: "Database with no port config",
			config: Config{
				Host: "localhost",
				Name: "dbname",
			},
			failExpected: true,
		},
		{
			name: "Database with wrong port config",
			config: Config{
				Host: "localhost",
				Port: -1,
			},
			failExpected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()

			if err != nil && !tc.failExpected {
				t.Errorf("config is valid, but an error has raised: %s", err)
			}

			if err == nil && tc.failExpected {
				t.Errorf("config is invalid but error hasn't raised")
			}
		})
	}
}
