/*
/ 一个对话一个aihelper
*/

package aihelper

import (
	"context"
	"sync"
)

var ctx = context.Background()

// AI助手管理器，管理用户-会话-AIHelper的映射关系
type AIHelperManager struct {
	helpers map[string]map[string]*AIHelper // map[用户账号（唯一）]map[会话ID]*AIHelper
	mtx     sync.RWMutex
}

// 获取或创建AIHelper
func (m *AIHelperManager) GetOrCreateAIHelper(userName string, sessionID string, modelType string, config map[string]string) (*AIHelper, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// 获取用户的会话映射
	userHelpers, exists := m.helpers[userName]
	if !exists {
		userHelpers = make(map[string]*AIHelper)
		m.helpers[userName] = userHelpers
	}

	helper, exists := userHelpers[sessionID]
	// 检查会话是否已存在
	if exists {
		return helper, nil
	} else {
		// 创建新的AIHelper
		helper, err := GetGlobalFactory().CreateAIHelper(ctx, modelType, sessionID, config)
		if err != nil {
			return nil, err
		}

		userHelpers[sessionID] = helper
		return helper, nil
	}
}

// 获取指定用户的指定会话的AIHelper
func (m *AIHelperManager) GetAIHelper(userName string, sessionID string) (*AIHelper, bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	userHelpers, exists := m.helpers[userName]
	if !exists {
		return nil, false
	}

	helper, exists := userHelpers[sessionID]
	return helper, exists
}

// 移除指定用户的指定会话的AIHelper
func (m *AIHelperManager) RemoveAIHelper(userName string, sessionID string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	userHelpers, exists := m.helpers[userName]
	if !exists {
		return
	}

	delete(userHelpers, sessionID)

	// 如果用户没有会话了，清理用户映射
	if len(userHelpers) == 0 {
		delete(m.helpers, userName)
	}
}

// 获取指定用户的所有会话ID
func (m *AIHelperManager) GetUserSessions(userName string) []string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	userHelpers, exists := m.helpers[userName]
	if !exists {
		return []string{}
	}

	sessionIDs := make([]string, 0, len(userHelpers))
	//取出所有的key
	for sessionID, _ := range userHelpers {
		sessionIDs = append(sessionIDs, sessionID)
	}

	return sessionIDs
}

// 全局管理器实例
var globalManager *AIHelperManager
var managerOnce sync.Once

// 创建新的管理器实例
func newAIHelperManager() *AIHelperManager {
	return &AIHelperManager{
		helpers: make(map[string]map[string]*AIHelper),
	}
}

func GetGlobalManager() *AIHelperManager {
	managerOnce.Do(func() {
		globalManager = newAIHelperManager()
	})
	return globalManager
}
