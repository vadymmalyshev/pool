package redis

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	red "github.com/go-redis/redis"
)

type RedisRepositorer interface {
	GetLatestWorker(walletId string) (map[string]string, error)
}

type RedisRepository struct {
	redisClient *red.Client
}

func NewRedisRepository(redisClient *red.Client) RedisRepositorer {
	return &RedisRepository{redisClient: redisClient}
}

func (repo *RedisRepository) GetLatestWorker(walletId string) (map[string]string, error) {
	result, err := repo.redisClient.HGetAll("last-update:" + walletId).Result()

	if apierrors.HandleError(err) {
		return make(map[string]string), err
	}
	return result, nil
}
