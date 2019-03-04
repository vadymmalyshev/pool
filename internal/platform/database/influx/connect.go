package influx

import (
	"fmt"
	"net/url"

	"git.tor.ph/hiveon/pool/internal/platform/database"
	client "github.com/influxdata/influxdb1-client"
	// . "github.com/influxdata/influxdb1-client/models"
)

// Connect returns initialized connection to db
func Connect(c database.Config) (*client.Client, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", c.Host, c.Port))
	if err != nil {
		return nil, err
	}

	client, err := client.NewClient(client.Config{URL: *u})
	if err != nil {
		return nil, err
	}

	client.SetAuth(c.User, c.Pass)
	if _, _, err := client.Ping(); err != nil {
		return nil, err
	}

	return client, nil
}
