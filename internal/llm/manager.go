package llm

import (
	"context"
	"fmt"
	"sync"
)

// AgentManager 智能体管理器 每个用户每个会话确定一个智能体
type AgentManager struct {
	lock    sync.RWMutex
	manager map[string]map[string]*Agent
}

var (
	globalAgentManager *AgentManager
	managerOnce        sync.Once //全局智能体管理器单例
)

// GetGlobalAgentManager 获取全局智能体管理器
func GetGlobalAgentManager() *AgentManager {
	managerOnce.Do(func() {
		globalAgentManager = &AgentManager{
			lock:    sync.RWMutex{},
			manager: make(map[string]map[string]*Agent),
		}
	})
	return globalAgentManager
}

// GetOrCreateAgent 获取或创建智能体
func (am *AgentManager) GetOrCreateAgent(ctx context.Context, username, sessionID string, modelType LLMType) (*Agent, error) {
	// 第一步：使用读锁检查智能体是否已经存在，减少锁竞争
	am.lock.RLock()
	userMap, userExists := am.manager[username]
	var existingAgent *Agent
	if userExists {
		existingAgent = userMap[sessionID]
	}
	am.lock.RUnlock()

	// 如果智能体已经存在且非空，直接返回
	if existingAgent != nil {
		return existingAgent, nil
	}

	// 第二步：使用写锁创建智能体
	am.lock.Lock()

	// double check：再次确认智能体是否已经被其他协程创建
	if userMap, userExists := am.manager[username]; userExists {
		if agent := userMap[sessionID]; agent != nil {
			return agent, nil
		}
	} else {
		// 用户名不存在，创建用户映射
		am.manager[username] = make(map[string]*Agent)
	}

	// 第三步：释放锁后再创建智能体，避免长时间持有锁
	// 先解锁
	am.lock.Unlock()

	// 创建智能体（耗时操作，不持有锁）
	factory := GetGlobalLLMFactory()
	newAgent, err := factory.CreateAgent(ctx, username, sessionID, modelType)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %v", err)
	}

	// 第四步：重新加锁，将创建好的智能体加入管理器
	am.lock.Lock()
	defer am.lock.Unlock()

	// 再次检查，确保没有其他协程已经创建了该智能体
	if existingAgent := am.manager[username][sessionID]; existingAgent != nil {
		return existingAgent, nil
	}

	// 将新创建的智能体加入管理器
	am.manager[username][sessionID] = newAgent
	return newAgent, nil
}

// GetAgent 获取指定用户指定会话的智能体
func (am *AgentManager) GetAgent(username, sessionID string) (*Agent, error) {
	am.lock.RLock()
	defer am.lock.RUnlock()
	if agent, ok := am.manager[username][sessionID]; ok {
		return agent, nil
	}
	return nil, fmt.Errorf("agent not found")
}

// RemoveAgent 移除指定用户指定会话的智能体
func (am *AgentManager) RemoveAgent(username, sessionID string) error {
	am.lock.Lock()
	defer am.lock.Unlock()
	if _, ok := am.manager[username]; !ok {
		return fmt.Errorf("user not found")
	}
	if _, ok := am.manager[username][sessionID]; !ok {
		return fmt.Errorf("agent not found")
	}
	// 移除智能体
	delete(am.manager[username], sessionID)
	// 移除用户
	if len(am.manager[username]) == 0 {
		delete(am.manager, username)
	}
	return nil
}

// GetUserSessionIDs 获取指定用户的所有会话ID
func (am *AgentManager) GetUserSessionIDs(username string) []string {
	am.lock.RLock()
	defer am.lock.RUnlock()
	sessionIDs := make([]string, 0)
	// 用户不存在或用户会话为空
	if userMap, ok := am.manager[username]; !ok || userMap == nil {
		return sessionIDs
	}
	userMap := am.manager[username]
	// 遍历用户会话，收集会话ID
	for sessionID := range userMap {
		sessionIDs = append(sessionIDs, sessionID)
	}
	return sessionIDs
}
