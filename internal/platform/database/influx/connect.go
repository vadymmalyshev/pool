package influx

import (
	"fmt"
	"net/url"

	client "github.com/influxdata/influxdb1-client"
)

// Connect returns initialized connection to db
func (db *DB) Connect() (*client.Client, error) {
	if err := db.Validate(); err != nil {
		return nil, err
	}
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", db.Host, db.Port))
	if err != nil {
		return nil, err
	}

	client, err := client.NewClient(client.Config{URL: *u})
	if err != nil {
		return nil, err
	}

	client.SetAuth(db.User, db.Pass)
	if _, _, err := client.Ping(); err != nil {
		return nil, err
	}

	return client, nil
}
