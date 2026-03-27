/*
/ factory 是生产aihelper的工厂
*/

package aihelper

import (
	"context"
	"fmt"
	"sync"
)

// ModelCreator 定义模型创建函数类型（需要 context）
type ModelCreator func(ctx context.Context, config map[string]string) (RagAIModel, error)

// AIModelFactory AI模型工厂
type AIModelFactory struct {
	creators map[string]ModelCreator
}

// 注册模型
func (f *AIModelFactory) registerCreators() {
	//DOUBAO_SEED_20
	f.creators[DOUBAO_SEED_20] = func(ctx context.Context, config map[string]string) (RagAIModel, error) {
		return NewRagArkModel(ctx, DOUBAO_SEED_20, config["username"], config["sessionId"])
	}
	//DEEPSEEK_V32
	f.creators[DEEPSEEK_V32] = func(ctx context.Context, config map[string]string) (RagAIModel, error) {
		return NewRagArkModel(ctx, DEEPSEEK_V32, config["username"], config["sessionId"])
	}
}

// 根据类型创建 AI 模型
func (f *AIModelFactory) createAIModel(ctx context.Context, modelType string, config map[string]string) (RagAIModel, error) {
	creator, ok := f.creators[modelType]
	if !ok {
		return nil, fmt.Errorf("unsupported model type: %s", modelType)
	}
	return creator(ctx, config)
}

// 一键创建 AIHelper
func (f *AIModelFactory) CreateAIHelper(ctx context.Context, modelType string, SessionID string, config map[string]string) (*AIHelper, error) {
	model, err := f.createAIModel(ctx, modelType, config)
	if err != nil {
		return nil, err
	}
	return NewAIHelper(model, SessionID), nil
}

// 可扩展注册
func (f *AIModelFactory) RegisterModel(modelType string, creator ModelCreator) {
	f.creators[modelType] = creator
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
