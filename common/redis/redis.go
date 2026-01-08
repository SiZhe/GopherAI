package redis

import (
	"GopherAI/config"
	"context"
	"strconv"
	"strings"
	"time"

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

// 设置验证码
func SetCaptchaForEmail(email, captcha string) error {
	// captcha:%s
	key := GenerateCaptcha(email)
	// 过期时间五分钟
	expire := 5 * time.Minute
	return Rdb.Set(ctx, key, captcha, expire).Err()
}

func CheckCaptchaForEmail(email, userInput string) (bool, error) {
	key := GenerateCaptcha(email)

	storedCaptcha, err := Rdb.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	// 忽略大小写的比较
	flag := strings.EqualFold(storedCaptcha, userInput)

	if flag {
		// 验证后删除
		Rdb.Del(ctx, key)
	}
	return flag, nil
}
