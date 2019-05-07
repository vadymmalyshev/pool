package redis

import (
	redisadapter "github.com/casbin/redis-adapter"
)

// Adapter returns adapter to redis
func (db *DB) Casbin() *redisadapter.Adapter {
	adaper := redisadapter.NewAdapter("tcp", db.Connection())
	return adaper
}
