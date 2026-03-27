package redis

import (
	"GopherAI/config"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var ctx = context.Background()

// 初始化redis
func InitRedis() {
	host := config.GetConfig().RedisConfig.RedisHost
	port := config.GetConfig().RedisConfig.RedisPort
	password := config.GetConfig().RedisConfig.RedisPassword
	db := config.GetConfig().RedisDb
	addr := host + ":" + strconv.Itoa(port)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
