package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Connect returns initialized connection to redis
func (db *DB) Connect() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", db.Connection())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//Connection returns connection string to redis server
func (db *DB) Connection() string {
	return fmt.Sprintf("%s:%d", db.Host, db.Port)
}
