package repository

import (
	"git.tor.ph/hiveon/pool/config"
	"strconv"
	red "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)
type IRedisRepository interface {
	GetLatestWorker(walletId string)  map[string]string
}

type RedisRepository struct {
	redisClient *red.Client
}

func NewRedisRepository() IRedisRepository {
	return &RedisRepository{redisClient: GetRedisClient()}
}

func GetRedisClient() *red.Client {
	DBName, _ := strconv.Atoi(config.Redis.Name)
	client := red.NewClient(&red.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: config.Redis.Pass,
		DB:       DBName,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Error(err)
	}

	return client
}

func (repo *RedisRepository) GetLatestWorker(walletId string) map[string]string {
	result, err := repo.redisClient.HGetAll("last-update:" + walletId).Result()

	if err != nil {
		log.Error(err)
	}
	return  result
}
