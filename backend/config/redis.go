package config

import (
	"log"

	"github.com/go-redis/redis"
)

func ConnectRedis() *redis.Client {
	opt, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(opt)
	return client
}
