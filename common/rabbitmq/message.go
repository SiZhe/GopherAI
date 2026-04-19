package rabbitmq

import (
	"GopherAI/common/mysql/dao/message"
	"GopherAI/common/mysql/model"
	"encoding/json"

	"github.com/streadway/amqp"
)

type MessageMQParam struct {
	SessionID string `json:"session_id"`
	UserName  string `json:"user_name"`
	Content   string `json:"content"`
	IsUser    bool   `json:"is_user"`
	Role      string `json:"role"`
	ModelType string `json:"model_type"`
}

func GenerateMessageMQParam(sessionID string, content string, userName string, isUser bool, role string, modelType string) ([]byte, error) {
	param := MessageMQParam{
		SessionID: sessionID,
		Content:   content,
		UserName:  userName,
		IsUser:    isUser,
		Role:      role,
		ModelType: modelType,
	}
	data, _ := json.Marshal(param)
	return data, nil
}

// 消费者异步插入到数据库中
func MQMessage(msg *amqp.Delivery) error {
	var param MessageMQParam
	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		return err
	}
	newMsg := &model.Message{
		SessionID: param.SessionID,
		UserName:  param.UserName,
		Content:   param.Content,
		IsUser:    param.IsUser,
		ModelType: param.ModelType,
		Role:      param.Role,
	}
	//消费者异步插入到数据库中
	_, err1 := message.CreateMessage(newMsg)
	if err1 != nil {
		return err
	}
	return nil
}
