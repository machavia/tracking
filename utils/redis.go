package utils

import (
	"../config"
	"github.com/go-redis/redis"
	"time"
)

var (
	// RedisClient is the connection handle
	// for the database
	RedisClient *redis.Client
)

func OpenDbConnexion() {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:        config.Config.RedisHost,
		Password:    "",
		DB:          0,
		PoolSize:    100,
		PoolTimeout: 30 * time.Second,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

}
