package user

import (
	"GopherAI/common/code"
	"GopherAI/common/email"
	"GopherAI/common/mysql/dao/device"
	"GopherAI/common/mysql/dao/user"
	"GopherAI/common/mysql/model"
	"GopherAI/common/redis"
	"GopherAI/config"
	"GopherAI/router/middleware/jwt"
	"GopherAI/utils"
	"log"
	"time"

	JWT "github.com/golang-jwt/jwt/v4"
)

func Register(userEmail, password, captcha string, deviceId string) (at string, rt string, c code.Code) {
	var ok bool
	var userInformation *model.User

	//1:先判断用户是否已经存在了
	if ok, _ := user.IsExistUserByEmail(userEmail); ok {
		return "", "", code.CodeUserExist
	}

	//2:从redis中验证验证码是否有效
	if ok, _ := redis.CheckCaptchaForEmail(userEmail, captcha); !ok {
		return "", "", code.CodeInvalidCaptcha
	}

	//3：生成11位的账号
	username := utils.GetRandomNumbers(11)

	//4：注册到数据库中
	if userInformation, ok = user.CreateUser(username, userEmail, password); !ok {
		return "", "", code.CodeServerBusy
	}

	//5：将账号一并发送到对应邮箱上去，后续需要账号登录
	if err := email.SendCaptcha(userEmail, username, email.UserNameMsg); err != nil {
		return "", "", code.CodeServerBusy
	}

	// 6:生成Token
	accessExpireDuration := time.Duration(config.GetConfig().JwtConfig.AccessExpire) * time.Minute
	refreshExpireDuration := time.Duration(config.GetConfig().JwtConfig.RefreshExpire) * time.Hour

	//3:返回一个Token
	accessToken, err := jwt.GenerateToken(userInformation.ID, userInformation.Username, deviceId, accessExpireDuration)
	if err != nil {
		return "", "", code.CodeServerBusy
	}
	refreshToken, err := jwt.GenerateToken(userInformation.ID, userInformation.Username, deviceId, refreshExpireDuration)
	if err != nil {
		return "", " ", code.CodeServerBusy
	}

	// 从黑名单中删除 没必要，因为发布时间不一样，不可能有相同的
	//err = redis.DeleteFromBlackList(accessToken)
	//if err != nil {
	//	return "", " ", code.CodeServerBusy
	//}
	//err = redis.DeleteFromBlackList(refreshToken)
	//if err != nil {
	//	return "", " ", code.CodeServerBusy
	//}

	return accessToken, refreshToken, code.CodeSuccess
}

func Login(username, password string, deviceId string) (at string, rt string, c code.Code) {
	var userInformation *model.User
	var ok bool

	//1:判断用户是否存在
	if ok, userInformation = user.IsExistUserByName(username); !ok {
		if ok, userInformation = user.IsExistUserByEmail(username); !ok {
			return "", "", code.CodeUserNotExist
		}
	}

	//2:判断用户是否密码账号正确
	if userInformation.Password != utils.MD5(password) {
		return "", "", code.CodeInvalidPassword
	}

	accessExpireDuration := time.Duration(config.GetConfig().JwtConfig.AccessExpire) * time.Minute
	refreshExpireDuration := time.Duration(config.GetConfig().JwtConfig.RefreshExpire) * time.Hour

	//3:返回一个Token
	accessToken, err := jwt.GenerateToken(userInformation.ID, userInformation.Username, deviceId, accessExpireDuration)
	if err != nil {
		return "", "", code.CodeServerBusy
	}
	refreshToken, err := jwt.GenerateToken(userInformation.ID, userInformation.Username, deviceId, refreshExpireDuration)
	if err != nil {
		return "", " ", code.CodeServerBusy
	}

	// 将改用户加入的数据库device中
	_, err = device.CreateDevice(username, deviceId, time.Now(), accessToken, refreshToken)
	if err != nil {
		log.Printf(err.Error())
		//return "", "", code.CodeServerBusy
	}

	// 从黑名单中删除 没必要，因为发布时间不一样，不可能有相同的
	//err = redis.DeleteFromBlackList(accessToken)
	//if err != nil {
	//	return "", " ", code.CodeServerBusy
	//}
	//err = redis.DeleteFromBlackList(refreshToken)
	//if err != nil {
	//	return "", " ", code.CodeServerBusy
	//}

	return accessToken, refreshToken, code.CodeSuccess
}

// 发验证码
func SendCaptcha(email_ string) code.Code {
	captcha := utils.GetRandomNumbers(6)

	//1:先存放到redis
	if err := redis.SetCaptchaForEmail(email_, captcha); err != nil {
		return code.CodeServerBusy
	}

	//2:再进行远程发送
	if err := email.SendCaptcha(email_, captcha, email.CodeMsg); err != nil {
		return code.CodeServerBusy
	}

	return code.CodeSuccess
}

// 直接从 refresh token 解析 username，无需传入
func RefreshAccessTokenWithToken(refreshAccessToken string, deviceId string) (string, code.Code) {
	rtClaims, err := jwt.ParseToken(refreshAccessToken, deviceId)
	if err != nil {
		// 解析失败时，尝试不验证地解析以获取 username 来删除设备
		tempClaims := struct {
			Username string `json:"username"`
			JWT.RegisteredClaims
		}{}
		_, _, _ = new(JWT.Parser).ParseUnverified(refreshAccessToken, &tempClaims)
		if tempClaims.Username != "" {
			_ = device.DeleteDevice(tempClaims.Username, deviceId)
		}
		return "", code.CodeInvalidRefreshToken
	}

	accessExpireDuration := time.Duration(config.GetConfig().JwtConfig.AccessExpire) * time.Minute
	newAccessToken, err := jwt.GenerateToken(rtClaims.ID, rtClaims.Username, deviceId, accessExpireDuration)
	if err != nil {
		return "", code.CodeServerBusy
	}
	return newAccessToken, code.CodeSuccess
}

// 传入username 以备设备管理需要
func RefreshAccessToken(username string, refreshAccessToken string, deviceId string) (string, code.Code) {
	rtClaims, err := jwt.ParseToken(refreshAccessToken, deviceId)
	if err != nil {
		// rt过期了 就在数据库device中删除
		err = device.DeleteDevice(username, deviceId)
		if err != nil {
			return "", code.CodeServerBusy
		}
		return "", code.CodeInvalidRefreshToken
	}

	accessExpireDuration := time.Duration(config.GetConfig().JwtConfig.AccessExpire) * time.Minute
	newAccessToken, err := jwt.GenerateToken(rtClaims.ID, rtClaims.Username, deviceId, accessExpireDuration)
	if err != nil {
		return "", code.CodeServerBusy
	}
	return newAccessToken, code.CodeSuccess
}
