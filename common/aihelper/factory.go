package aihelper

import (
	"context"
	"fmt"
	"sync"
)

// ModelCreator 定义模型创建函数类型（需要 context）
type ModelCreator func(ctx context.Context, config map[string]interface{}) (AIModel, error)

// AIModelFactory AI模型工厂
type AIModelFactory struct {
	creators map[string]ModelCreator
}

// 注册模型
func (f *AIModelFactory) registerCreators() {
	//DOUBAO_SEED_V16
	f.creators[DOUBAO_SEED_V16] = func(ctx context.Context, config map[string]interface{}) (AIModel, error) {
		return NewARKModel(ctx, DOUBAO_SEED_V16)
	}
	//DOUBAO_SEED_V16_LITE
	f.creators[DOUBAO_SEED_V16_LITE] = func(ctx context.Context, config map[string]interface{}) (AIModel, error) {
		return NewARKModel(ctx, DOUBAO_SEED_V16_LITE)
	}
	//DOUBAO_SEED_CODE
	f.creators[DOUBAO_SEED_CODE] = func(ctx context.Context, config map[string]interface{}) (AIModel, error) {
		return NewARKModel(ctx, DOUBAO_SEED_CODE)
	}
	//DEEPSEEK_V32
	f.creators[DEEPSEEK_V32] = func(ctx context.Context, config map[string]interface{}) (AIModel, error) {
		return NewARKModel(ctx, DEEPSEEK_V32)
	}
}

// 全局管理器实例
var globalFactory *AIModelFactory
var factoryOnce sync.Once

// GetGlobalFactory 获取全局单例
func GetGlobalFactory() *AIModelFactory {
	factoryOnce.Do(func() {
		globalFactory = &AIModelFactory{
			creators: make(map[string]ModelCreator),
		}
		globalFactory.registerCreators()
	})
	return globalFactory
}

// CreateAIModel 根据类型创建 AI 模型
func (f *AIModelFactory) CreateAIModel(ctx context.Context, modelType string, config map[string]interface{}) (AIModel, error) {
	creator, ok := f.creators[modelType]
	if !ok {
		return nil, fmt.Errorf("unsupported model type: %s", modelType)
	}
	return creator(ctx, config)
}

// CreateAIHelper 一键创建 AIHelper
func (f *AIModelFactory) CreateAIHelper(ctx context.Context, modelType string, SessionID string, config map[string]interface{}) (*AIHelper, error) {
	model, err := f.CreateAIModel(ctx, modelType, config)
	if err != nil {
		return nil, err
	}
	return NewAIHelper(model, SessionID), nil
}

// RegisterModel 可扩展注册
func (f *AIModelFactory) RegisterModel(modelType string, creator ModelCreator) {
	f.creators[modelType] = creator
}
