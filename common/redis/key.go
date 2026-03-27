/*
/ redis的作用的 检查验证码是否有效（5分钟）
*/

package redis

import (
	"GopherAI/config"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// key:特定邮箱-> 验证码
func GenerateCaptcha(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
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
