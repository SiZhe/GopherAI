package session

import (
	"GopherAI/common/mysql"
	"GopherAI/common/mysql/model"
)

func CreateSession(session *model.Session) (*model.Session, error) {
	err := mysql.DB.Create(session).Error
	return session, err
}

func GetSessionsByUserName(UserName int64) ([]model.Session, error) {
	var sessions []model.Session
	err := mysql.DB.Where("user_name = ?", UserName).Find(&sessions).Error
	return sessions, err
}

func GetSessionByID(sessionID string) (*model.Session, error) {
	var session model.Session
	err := mysql.DB.Where("id = ?", sessionID).First(&session).Error
	return &session, err
}
