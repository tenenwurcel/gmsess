package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisCli *redis.Client

func SetupRedis() {
	host := os.Getenv("REDISHOST")
	port := os.Getenv("REDISPORT")

	RedisCli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: os.Getenv("REDISPWD"),
		DB:       0, // use default DB
	})

	setupRedisCli(RedisCli)
}

func setupRedisCli(RedisCli *redis.Client) {
	redisCli = RedisCli
}

func GetRedisCli() *redis.Client {
	return redisCli
}
