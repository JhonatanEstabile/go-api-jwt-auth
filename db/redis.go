package db

import (
	"os"

	"github.com/go-redis/redis"
)

var client *redis.Client

//GetRedisClient initialize a redis client
func GetRedisClient() *redis.Client {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}
