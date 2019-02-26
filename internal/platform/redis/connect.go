package redis

import "github.com/gomodule/redigo/redis"

// Connect returns initialized connection to redis
func Connect(c Config) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", c.Connection())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
