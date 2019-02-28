package redis

import (
	"fmt"

	"git.tor.ph/hiveon/pool/internal/platform/database"
	"github.com/gomodule/redigo/redis"
)

// Connect returns initialized connection to redis
func Connect(c database.Config) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", Connection(c))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//Connection returns connection string to redis server
func Connection(c database.Config) string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
