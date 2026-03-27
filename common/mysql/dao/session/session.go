package session

import (
	"GopherAI/common/mysql"
	"GopherAI/common/mysql/model"
	"GopherAI/rag"
	"fmt"
)

func CreateSession(session *model.Session) (*model.Session, error) {
	// 创建对话的时候创建一个文件夹
	err := rag.CreateUploadsDir(session.UserName, session.ID)
	fmt.Println(session.ID)
	if err != nil {
		return nil, err
	}

	err = mysql.DB.Create(session).Error
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
