package redis

import (
	"git.tor.ph/hiveon/pool/internal/platform/database"
	redisadapter "github.com/casbin/redis-adapter"
)

// Adapter returns adapter to redis
func Adapter(c database.Config) *redisadapter.Adapter {
	adaper := redisadapter.NewAdapter("tcp", Connection(c))

	return adaper
}
