package jwt

import (
	"GopherAI/common/redis"
	"GopherAI/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	DeviceId string `json:"deviceId"`
	jwt.RegisteredClaims
}

func GenerateToken(id int64, username string, deviceId string, expire time.Duration) (string, error) {
	claims := Claims{
		ID:       id,
		Username: username,
		DeviceId: deviceId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  config.GetConfig().Issuer,
			Subject: config.GetConfig().Subject,
			// 一小时的有效期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 加密
	return token.SignedString([]byte(config.GetConfig().Key))
}

// ParseToken 解析Token
// 1.验证是否在黑名单中 ; 2.解析是否有效 ; 3.设备id是否匹配，以免背窃取
func ParseToken(token string, deviceId string) (*Claims, error) {
	// 检查是否在黑名单中
	if redis.CheckTokenInBlackList(token) {
		return nil, jwt.ErrSignatureInvalid
	}

	// 检查是否有效
	claims := new(Claims)
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Key), nil
	})
	if err != nil {
		return nil, err
	}

	if !t.Valid || claims == nil {
		return nil, jwt.ErrSignatureInvalid
	}

	if claims.DeviceId != deviceId {
		return nil, errors.New("deviceId incompatible")
	}

	return claims, nil
}

// 得到token的过期时间戳
func GetTokenExpireTime(tokenStr string) (time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &jwt.RegisteredClaims{})
	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid claims")
	}

	return claims.ExpiresAt.Time, nil
}
