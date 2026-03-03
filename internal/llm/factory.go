package llm

import (
	"context"
	"fmt"
	"sync"
)

type LLMCreator func(ctx context.Context) (LLM, error)

// LLMFactory 模型工厂 每个模型类型对应一个模型工厂
type LLMFactory struct {
	// 模型工厂
	creators map[LLMType]LLMCreator
}

var (
	globalFactory *LLMFactory
	factoryOnce   sync.Once //全局工厂单例
)

// GetGlobalLLMFactory 获取全局模型工厂
func GetGlobalLLMFactory() *LLMFactory {
	factoryOnce.Do(func() {
		globalFactory = &LLMFactory{
			creators: make(map[LLMType]LLMCreator),
		}
		// 注册模型创建器
		globalFactory.registerCreators()
	})
	return globalFactory
}

// registerCreator 注册模型创建器 为每个类型注册一个模型创建器
func (f *LLMFactory) registerCreators() {
	// 注册火山引擎文本模型创建器
	f.creators[ArkModelType] = func(ctx context.Context) (LLM, error) {
		return NewArkTextModel(ctx)
	}
	// TODO 注册其他模型创建器 (目前只有火山引擎文本模型后续接入其他Eino支持的模型厂商)
}

// CreateLLM 创建模型
func (f *LLMFactory) CreateLLM(ctx context.Context, llmType LLMType) (LLM, error) {
	creator, ok := f.creators[llmType]
	if !ok {
		return nil, fmt.Errorf("no creator found for LLM type: %v", llmType)
	}
	return creator(ctx)
}

// CreateAgent 创建智能体
func (f *LLMFactory) CreateAgent(ctx context.Context, username, sessionID string, llmType LLMType) (*Agent, error) {
	llm, err := f.CreateLLM(ctx, llmType)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %v", err)
	}
	return NewAgent(username, sessionID, llm), nil
}

// RegisterLLM 注册模型创建器 自定义模型类型
func (f *LLMFactory) RegisterLLM(llmType string, creator LLMCreator) {
	f.creators[LLMType(llmType)] = creator
}
