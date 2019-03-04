package repository

// import (
// 	"git.tor.ph/hiveon/pool/config"
// 	"git.tor.ph/hiveon/pool/internal/platform/database/redis"
// 	red "github.com/gomodule/redigo/redis"
// 	log "github.com/sirupsen/logrus"
// )

// type IRedisRepository interface {
// 	GetLatestWorker(walletId string) map[string]string
// }

// type RedisRepository struct {
// 	redisClient red.Conn
// }

// func NewRedisRepository() IRedisRepository {
// 	client, err := redis.Connect(config.Redis)

// 	if err != nil {
// 		log.Panic("failed to init redis db :", err.Error())
// 	}

// 	return &RedisRepository{redisClient: client}
// }

// func (repo *RedisRepository) GetLatestWorker(walletId string) map[string]string {
// 	//result, err := repo.redisClient.HGetAll("last-update:" + walletId).Result()
// 	result, err := repo.redisClient.Do("HGETALL", "last-update:"+walletId)

// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return result.(map[string]string)
// }
