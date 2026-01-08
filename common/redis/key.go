package redis

import (
	"GopherAI/config"
	"fmt"
)

// key:特定邮箱-> 验证码
func GenerateCaptcha(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
}
