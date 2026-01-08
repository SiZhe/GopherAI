package aihelper

import (
	"GopherAI/common/mysql/model"
	"GopherAI/common/rabbitmq"
	"GopherAI/utils"
	"context"
	"sync"
)

type AIHelper struct {
	model    AIModel
	messages []*model.Message
	mtx      sync.RWMutex
	//一个会话绑定一个AIHelper
	SessionID string
	saveFunc  func(*model.Message) (*model.Message, error)
}

// NewAIHelper 创建新的AIHelper实例
func NewAIHelper(model_ AIModel, SessionID string) *AIHelper {
	return &AIHelper{
		model:     model_,
		messages:  make([]*model.Message, 0),
		SessionID: SessionID,
		//异步推送到消息队列中
		saveFunc: func(msg *model.Message) (*model.Message, error) {
			data, err := rabbitmq.GenerateMessageMQParam(msg.SessionID, msg.Content, msg.UserName, msg.IsUser, msg.ModelType)
			if err != nil {
				return msg, err
			}
			err = rabbitmq.RMQMessage.Publish(data)
			return msg, err
		},
	}
}

// addMessage 添加消息到内存中并调用自定义存储函数(存储到rabbitmq中)
func (a *AIHelper) AddMessage(content string, userName string, isUser bool, modelType string, save bool) {
	msg := model.Message{
		SessionID: a.SessionID,
		UserName:  userName,
		Content:   content,
		IsUser:    isUser,
		ModelType: modelType,
	}
	a.messages = append(a.messages, &msg)
	if save {
		_, err := a.saveFunc(&msg)
		if err != nil {
			return
		}
	}
}

// SaveMessage 保存消息到数据库（通过回调函数避免循环依赖）
// 通过传入func，自己调用外部的保存函数，即可支持同步异步等多种策略
func (a *AIHelper) SetSaveFunc(saveFunc func(*model.Message) (*model.Message, error)) {
	a.saveFunc = saveFunc
}

// GetMessages 获取所有消息历史
func (a *AIHelper) GetMessages() []*model.Message {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	out := make([]*model.Message, len(a.messages))
	copy(out, a.messages)
	return out
}

// GetModelType 获取模型类型
func (a *AIHelper) GetModelType() string {
	return a.model.GetModelType()
}

// 同步生成
func (a *AIHelper) GenerateResponse(userName string, ctx context.Context, userQuestion string) (*model.Message, error) {
	//调用存储函数
	a.AddMessage(userQuestion, userName, true, a.model.GetModelType(), true)

	a.mtx.RLock()
	//将model.Message转化成schema.Message
	messages := utils.ConvertToSchemaMessages(a.messages)
	a.mtx.RUnlock()

	//调用模型生成回复
	schemaMsg, err := a.model.GenerateResponse(ctx, messages)
	if err != nil {
		return nil, err
	}

	//将schema.Message转化成model.Message
	modelMsg := utils.ConvertToModelMessage(a.SessionID, userName, schemaMsg)

	//调用存储函数
	a.AddMessage(modelMsg.Content, userName, false, a.model.GetModelType(), true)

	return modelMsg, nil
}

// 流式生成
func (a *AIHelper) StreamResponse(userName string, ctx context.Context, cb StreamCallback, userQuestion string) (*model.Message, error) {
	//调用存储函数
	a.AddMessage(userQuestion, userName, true, a.model.GetModelType(), true)

	a.mtx.RLock()
	messages := utils.ConvertToSchemaMessages(a.messages)
	a.mtx.RUnlock()

	content, err := a.model.StreamResponse(ctx, messages, cb)
	if err != nil {
		return nil, err
	}
	//转化成model.Message
	modelMsg := &model.Message{
		SessionID: a.SessionID,
		UserName:  userName,
		Content:   content,
		IsUser:    false,
	}

	//调用存储函数
	a.AddMessage(modelMsg.Content, userName, false, a.model.GetModelType(), true)

	return modelMsg, nil
}
