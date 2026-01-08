package user

import (
	"GopherAI/common/code"
	"GopherAI/common/email"
	"GopherAI/common/mysql/dao/user"
	"GopherAI/common/mysql/model"
	"GopherAI/common/redis"
	"GopherAI/router/middleware/jwt"
	"GopherAI/utils"
)

func Register(userEmail, password, captcha string) (string, code.Code) {
	var ok bool
	var userInformation *model.User

	//1:先判断用户是否已经存在了
	if ok, _ := user.IsExistUserByEmail(userEmail); ok {
		return "", code.CodeUserExist
	}

	//2:从redis中验证验证码是否有效
	if ok, _ := redis.CheckCaptchaForEmail(userEmail, captcha); !ok {
		return "", code.CodeInvalidCaptcha
	}

	//3：生成11位的账号
	username := utils.GetRandomNumbers(11)

	//4：注册到数据库中x
	if userInformation, ok = user.CreateUser(username, userEmail, password); !ok {
		return "", code.CodeServerBusy
	}

	//5：将账号一并发送到对应邮箱上去，后续需要账号登录
	if err := email.SendCaptcha(userEmail, username, email.UserNameMsg); err != nil {
		return "", code.CodeServerBusy
	}

	// 6:生成Token
	token, err := jwt.GenerateToken(userInformation.ID, userInformation.Username)

	if err != nil {
		return "", code.CodeServerBusy
	}

	return token, code.CodeSuccess
}

func Login(countName, password string) (string, code.Code) {
	var userInformation *model.User
	var ok bool
	//1:判断用户是否存在

	if ok, userInformation = user.IsExistUserByName(countName); !ok {
		if ok, userInformation = user.IsExistUserByEmail(countName); !ok {
			return "", code.CodeUserNotExist
		}
	}

	//2:判断用户是否密码账号正确
	if userInformation.Password != utils.MD5(password) {
		return "", code.CodeInvalidPassword
	}
	//3:返回一个Token
	token, err := jwt.GenerateToken(userInformation.ID, userInformation.Username)

	if err != nil {
		return "", code.CodeServerBusy
	}
	return token, code.CodeSuccess
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
