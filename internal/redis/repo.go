package redis

import (
	red "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type RedisRepositorer interface {
	GetLatestWorker(walletId string)  map[string]string
}

type RedisRepository struct {
	redisClient *red.Client
}

func NewRedisRepository(redisClient *red.Client) RedisRepositorer {
	return &RedisRepository{redisClient: redisClient}
}

func (repo *RedisRepository) GetLatestWorker(walletId string) map[string]string {
	result, err := repo.redisClient.HGetAll("last-update:" + walletId).Result()

	if err != nil {
		log.Error(err)
	}
	return  result
}
