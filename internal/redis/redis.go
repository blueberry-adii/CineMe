package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

/*
* Constructor which returns a new Redis Client
* required by Redis Store
 */
func NewRedisClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}

	log.Printf("connected to redis at %s", addr)

	return rdb
}
