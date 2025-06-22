package configs

import (
	"time"

	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/redis"
)

var RedisClient *redis.Client

func InitRedis() {
	log.Infof("Initializing Redis...")
	config := &redis.Config{
		Hosts: []string{"localhost:6379"},
		PoolSize: 50,
		MinIdleConn: 10,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  10 * time.Minute,
	}
	
	RedisClient = redis.NewClient(config)
	log.Infof("Redis initialized successfully!")
}
