package user

import (
	"GopherAI/common/mysql"
	"GopherAI/common/mysql/model"
	"GopherAI/utils"

	"gorm.io/gorm"
)

func IsExistUserByName(username string) (bool, *model.User) {
	var user model.User
	err := mysql.DB.Where("username = ?", username).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return true, &user
}

func IsExistUserByEmail(email string) (bool, *model.User) {
	var user model.User
	err := mysql.DB.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return true, &user
}

func CreateUser(username string, email string, password string) (*model.User, bool) {
	user := &model.User{
		Name:     username,
		Email:    email,
		Username: username,
		Password: utils.MD5(password),
	}

	err := mysql.DB.Create(&user).Error

	if err != nil {
		return nil, false
	}

	return user, true
}

func GetUsernameByEmail(email string) (error, string) {
	var user model.User
	err := mysql.DB.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return err, ""
	}

	return nil, user.Username
}
