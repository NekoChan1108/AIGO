package llm

import (
	"AIGO/internal/dao"
	"AIGO/internal/model"
	"AIGO/pkg/log"
	"AIGO/pkg/mq"
	"AIGO/pkg/utils/message"
	"AIGO/pkg/utils/prompt"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// TODO  解决文件上传刷新页面导致附加丢失不显示的问题

// Agent 智能会话
type Agent struct {
	lock      sync.RWMutex // 读写锁 保护上下文消息的并发访问
	model     LLM          // Agent对应的模型
	username  string       // 用户名
	sessionID string       // 当前的会话ID
	// contextMsg []*model.Message                                        // 上下文消息 当前会话下的上下文消息
	saveFunc func(ctx context.Context, message *model.Message) error // 保存会话记录 持久化当前会话的消息到数据库
}

// NewAgent 创建智能会话
func NewAgent(username, sessionID string, llm LLM) *Agent {
	return &Agent{
		lock: sync.RWMutex{},
		// contextMsg: make([]*model.Message, 0),
		sessionID: sessionID,
		username:  username,
		model:     llm,
		saveFunc: func(ctx context.Context, message *model.Message) error {
			data, err := json.Marshal(message)
			if err != nil {
				return fmt.Errorf("marshal message failed: %w", err)
			}
			// 写入kafka
			err = mq.KafkaProducer.WriteMessages(ctx, kafka.Message{
				Key:   []byte(message.ID), // 使用消息ID作为key保证消息的唯一性
				Value: data,
			})
			if err != nil {
				return fmt.Errorf("write message failed: %w", err)
			}
			return nil
		},
	}
}

// addMessage 添加消息
func (a *Agent) addMessage(ctx context.Context, id, content string, isUser, save bool, createdAt time.Time) (*model.Message, error) {
	// 构建消息
	message := &model.Message{
		ID:        id,
		SessionId: a.sessionID,
		Content:   content,
		Username:  a.username,
		IsUser:    isUser,
		CreatedAt: createdAt,
	}
	// 添加到上下文消息
	// a.contextMsg = append(a.contextMsg, message)
	// 异步保存消息到数据库 同时将消息返回给调用者防止出现消息丢失的情况
	if save {
		return message, a.saveFunc(ctx, message)
	}
	return message, nil
}

// SetSaveFunc 设置保存会话记录的函数 方便后续扩展
func (a *Agent) SetSaveFunc(saveFunc func(ctx context.Context, message *model.Message) error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.saveFunc = saveFunc
}

// GetContextMsg 获取上下文消息
func (a *Agent) getContextMsg(ctx context.Context) ([]*model.Message, error) {
	messages, err := dao.GetMessagesBySessionID(ctx, a.sessionID)
	if err != nil {
		return nil, fmt.Errorf("get messages by sessionID failed: %w", err)
	}
	return messages, nil
}

// GenerateResponse 生成响应
func (a *Agent) GenerateResponse(ctx context.Context, question string) (*model.Message, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	// 先将问题添加到上下文消息
	questionID := uuid.New().String()
	userMsg, err := a.addMessage(ctx, questionID, question, true, true, time.Now())
	if err != nil {
		return nil, fmt.Errorf("add message failed: %w", err)
	}
	messages := make([]*model.Message, 0)
	// 从mongo中读取上下文消息
	ctxMsg, err := a.getContextMsg(ctx)
	if err != nil {
		return nil, fmt.Errorf("get context messages failed: %w", err)
	}
	// 确保上下文消息中最后一条消息是用户发送的 否则说明上下文消息中缺少用户发送的消息 保证上下文消息的完整性
	if len(ctxMsg) > 0 && ctxMsg[len(ctxMsg)-1].IsUser && ctxMsg[len(ctxMsg)-1].ID == userMsg.ID {
		log.Infof("context messages is complete: %v", ctxMsg)
		messages = ctxMsg
	} else {
		// 添加用户发送的问题 防止异步写入导致读取不到用户发送的问题
		messages = append(ctxMsg, userMsg)
	}
	// 将消息转为AI模型接受的消息
	schemaMessages := message.ConvertToSchemaMessages(messages)
	// 调用AI模型生成响应
	response, err := a.model.GenerateResponse(ctx, schemaMessages)
	if err != nil {
		return nil, fmt.Errorf("generate response failed: %w", err)
	}
	// 将结果转为消息存储到数据库
	message := message.ConvertToModelMessage(a.username, a.sessionID, response)
	// 添加到上下文消息
	a.addMessage(ctx, message.ID, message.Content, false, true, message.CreatedAt)
	return message, nil
}

// StreamResponse 流式响应
func (a *Agent) StreamResponse(ctx context.Context, question string, callback StreamCallback) (*model.Message, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	// 先将问题添加到上下文消息
	questionID := uuid.New().String()
	userMsg, err := a.addMessage(ctx, questionID, question, true, true, time.Now())
	if err != nil {
		return nil, fmt.Errorf("add message failed: %w", err)
	}
	messages := make([]*model.Message, 0)
	// 从mongo中读取上下文消息
	ctxMsg, err := a.getContextMsg(ctx)
	if err != nil {
		return nil, fmt.Errorf("get context messages failed: %w", err)
	}
	if len(ctxMsg) > 0 && ctxMsg[len(ctxMsg)-1].IsUser && ctxMsg[len(ctxMsg)-1].ID == userMsg.ID {
		log.Infof("context messages is complete: %v", ctxMsg)
		messages = ctxMsg
	} else {
		// 添加用户发送的问题 防止异步写入导致读取不到用户发送的问题
		messages = append(ctxMsg, userMsg)
	}
	schemaMessages := message.ConvertToSchemaMessages(messages)
	// 获取RAG检索器
	retriever, err := dao.GetRagRetriever(ctx, a.username, a.sessionID)
	if err != nil {
		return nil, fmt.Errorf("get rag retriever failed: %w", err)
	}
	// 如果存在说明用户上传了文件 优先使用RAG检索器
	if retriever != nil {
		// 召回 TOPK 条文档
		docs, err := retriever.RetrieveFile(ctx, question)
		if err != nil {
			return nil, fmt.Errorf("retrieve file failed: %w", err)
		}
		if len(docs) > 0 {
			// 构建提示词
			prompt := prompt.BuildPrompt(question, docs)
			if len(schemaMessages) > 0 {
				// 替换最后一条消息的内容为构建好的提示词
				lastMsg := schemaMessages[len(schemaMessages)-1]
				lastMsg.Content = prompt
			}
		}
	}
	// 调用AI模型生成响应
	response, err := a.model.StreamResponse(ctx, schemaMessages, callback)
	if err != nil {
		return nil, fmt.Errorf("stream response failed: %w", err)
	}
	msg := &model.Message{
		ID:        uuid.New().String(),
		SessionId: a.sessionID,
		Content:   response,
		Username:  a.username,
		IsUser:    false,
		CreatedAt: time.Now(),
	}
	// 添加到上下文消息
	if _, err := a.addMessage(ctx, msg.ID, response, false, true, msg.CreatedAt); err != nil {
		return nil, fmt.Errorf("add message failed: %w", err)
	}
	return msg, nil
}

// GetModelType 获取模型类型
func (a *Agent) GetModelType() LLMType {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.model.GetModelType()
}
