package redis

import (
	"git.tor.ph/hiveon/pool/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"git.tor.ph/hiveon/pool/internal/consumer/utils"
	"strconv"
	"strings"
	"time"
)
type IRedisRepository interface {
	RedisCount(name string)
	RedisSet(wallet string, param map[string]interface{})
	RedisDel(wallet string, key string)
	RedisAlive()(s string)
	RedisGetPoints(hashId string) (map[string]string)
}

type RedisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository() IRedisRepository {
	return &RedisRepository{redisClient: GetRedisClient()}
}

func GetRedisClient() *redis.Client {
	host := config.Redis.Host
	port := config.Redis.Port
	password := config.Redis.Pass
	db, _ := strconv.Atoi(config.Redis.Name)
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Error(err)
	}

	return client
}

func (repo *RedisRepository) RedisCount(name string) {
	now := time.Now()
	f := []string {strconv.Itoa(now.Hour()),":",strconv.Itoa(now.Minute())}
	field := strings.Join(f,"")
	k := []string {strconv.Itoa(int(now.Month())),"-",strconv.Itoa(now.Day())}
	key := strings.Join(k,"")
	keyArray :=[]string {"count:", name, ":", key}
	hashId := strings.Join(keyArray,"")
	repo.redisClient.HIncrBy(hashId, field,1)
}

func (repo *RedisRepository) RedisSet(wallet string, param map[string]interface{}) {
	strArray := []string{"last-update:",wallet}
	key := strings.Join(strArray,"")
	repo.redisClient.HMSet(key, param)
}

func (repo *RedisRepository) RedisDel(wallet string, value string) {
	strArray := []string{"offline-notified:",wallet}
	key := strings.Join(strArray,"")
	repo.redisClient.HDel(key, value)
}

func (repo *RedisRepository) RedisAlive() (s string) {
	_, err := repo.redisClient.Ping().Result()
	if err != nil {
		return utils.IsDown
	}
	return utils.IsUP
}

func (repo *RedisRepository) RedisGetPoints(hashId string) (map[string]string) {
	return repo.redisClient.HGetAll(hashId).Val()
}

