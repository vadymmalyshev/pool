package redis

import (
	"fmt"
	"strconv"
	"github.com/go-redis/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// Connect returns initialized connection to redis
func (db *DB) Connect() (*redis.Client, error) {

	host := db.Host
	port := db.Port
	password := db.Pass
	dbname, _ := strconv.Atoi(db.Name)

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: password,
		DB:       dbname,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (db *DB) ConnectCasbin() (redigo.Conn, error) {
	conn, err := redigo.Dial("tcp", db.Connection())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//Connection returns connection string to redis server
func (db *DB) Connection() string {
	return fmt.Sprintf("%s:%d", db.Host, db.Port)
}
