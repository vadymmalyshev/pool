package redis

import redisadapter "github.com/casbin/redis-adapter"

// Adapter returns adapter to redis
func Adapter(c Config) *redisadapter.Adapter {
	adaper := redisadapter.NewAdapter("tcp", c.Connection())

	return adaper
}
