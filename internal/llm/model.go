package llm

import (
	"AIGO/config"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

// LLMType 模型类型 后续拓展Eino支持的不同厂商的模型例如OpenAI Qwen QianFan 等
type LLMType string

const (
	// ArkModelType 火山模型
	ArkModelType LLMType = "ark"
	// OpenAIModelType OpenAI模型
	OpenAIModelType LLMType = "openai"
	// QianFanModelType 千帆模型
	QianFanModelType LLMType = "qianfan"
	// QwenModelType Qwen模型
	QwenModelType LLMType = "qwen"
	// ArkBotModelType 火山引擎机器人模型
	ArkBotModelType LLMType = "ark_bot"
	// DeepSeekModelType DeepSeek模型
	DeepSeekModelType LLMType = "deepseek"
	// GeminiModelType Gemini模型
	GeminiModelType LLMType = "gemini"
	// ClaudeModelType Claude模型
	ClaudeModelType LLMType = "claude"
	// OllamaModelType Ollama模型
	OllamaModelType LLMType = "ollama"
)

// StreamCallback 流式回调函数
type StreamCallback func(string) error

// LLM 大模型接口 方便后续不同厂商的模型实现
type LLM interface {
	// GenerateResponse 生成响应
	GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
	// StreamResponse 流式响应
	StreamResponse(ctx context.Context, messages []*schema.Message, callback StreamCallback) (string, error)
	// GetModelType 获取模型类型
	GetModelType() LLMType
}

// ArkTextModel 火山引擎文本模型
type ArkTextModel struct {
	llm *ark.ChatModel
}

// NewArkTextModel 创建火山引擎文本模型
func NewArkTextModel(ctx context.Context) (*ArkTextModel, error) {
	llm, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		Model:  config.Cfg.ModelCfg.TextModel,
		APIKey: config.Cfg.ModelCfg.ApiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Ark model: %v", err)
	}
	return &ArkTextModel{llm: llm}, nil
}

// GenerateResponse 生成响应
func (m *ArkTextModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	msg, err := m.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %v", err)
	}
	return msg, nil
}

// StreamResponse 流式响应
func (m *ArkTextModel) StreamResponse(ctx context.Context, messages []*schema.Message, callback StreamCallback) (string, error) {
	reader, err := m.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to stream response: %v", err)
	}
	defer reader.Close()
	// 初始化防止空指针
	fullContent := &strings.Builder{}
	for {
		msg, err := reader.Recv()
		// 读取到EOF时，跳出循环
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read stream response: %v", err)
		}
		// 处理msg 先写入fullContent
		if len(msg.Content) > 0 {
			_, err := fullContent.WriteString(msg.Content)
			if err != nil {
				return "", fmt.Errorf("failed to write stream response: %v", err)
			}
			// 回调函数处理内容 流式传到前端
			err = callback(msg.Content)
			if err != nil {
				return "", fmt.Errorf("callback failed: %v", err)
			}
			fmt.Println(msg.Content)
		}
	}
	// 全部返回
	return fullContent.String(), nil
}

// GetModelType 获取模型类型
func (m *ArkTextModel) GetModelType() LLMType {
	return ArkModelType
}
