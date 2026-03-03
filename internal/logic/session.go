package logic

import (
	"AIGO/internal/dao"
	"AIGO/internal/llm"
	"AIGO/internal/model"
	"AIGO/pkg/log"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// GetSessionByUsername 获取用户的所有会话
func GetSessionByUsername(ctx context.Context, username string) ([]model.SessionInfo, error) {
	sessions, err := dao.GetSessionByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get sessions by username failed: %v", err)
	}
	sessionInfos := make([]model.SessionInfo, 0, len(sessions))
	for _, session := range sessions {
		sessionInfos = append(sessionInfos, model.SessionInfo{
			ID:    session.ID,
			Title: session.Title,
		})
	}
	return sessionInfos, nil
}

// DeleteSession 删除会话
func DeleteSession(ctx context.Context, sessionID string) error {
	res, err := dao.DeleteSession(ctx, sessionID)
	if err != nil || res <= 0 {
		return fmt.Errorf("delete session failed: %v", err)
	}
	// TODO 删除会话对应的agent 删除会话对应的文件以及 milvus 索引
	return nil
}

// CreateSessionAndSendMessage 创建会话并发送第一条消息 第一条消息必然是用户的第一条消息
func CreateSessionAndSendMessage(ctx context.Context, userQuestion, username string, llmType llm.LLMType) (string, string, error) {
	session := &model.Session{
		ID:        uuid.New().String(),
		Username:  username,
		CreatedAt: time.Now(),
		Title:     userQuestion,
	}
	// 写入数据库
	res, err := dao.CreateSession(ctx, session)
	if err != nil || res <= 0 {
		return "", "", fmt.Errorf("create session failed: %v", err)
	}
	// 为会话分配agent
	manager := llm.GetGlobalAgentManager()
	agent, err := manager.GetOrCreateAgent(ctx, username, session.ID, llmType)
	if err != nil {
		return "", "", fmt.Errorf("create agent failed: %v", err)
	}
	// 发送第一条消息
	msg, err := agent.GenerateResponse(ctx, userQuestion)
	if err != nil {
		return "", "", fmt.Errorf("generate response failed: %v", err)
	}
	return session.ID, msg.Content, nil
}

// ChatSend 聊天 会话中返回对问题的解答
func ChatSend(ctx context.Context, userQuestion, username, sessionID string, llmType llm.LLMType) (string, error) {
	manager := llm.GetGlobalAgentManager()
	agent, err := manager.GetOrCreateAgent(ctx, username, sessionID, llmType)
	if err != nil {
		return "", fmt.Errorf("get agent failed: %v", err)
	}
	msg, err := agent.GenerateResponse(ctx, userQuestion)
	if err != nil {
		return "", fmt.Errorf("generate response failed: %v", err)
	}
	return msg.Content, nil
}

// CreateStreamSession 创建流式会话 流式会话不发送第一条消息 仅创建会话
func CreateStreamSession(ctx context.Context, userQuestion string, username string) (string, error) {
	now := time.Now()
	session := &model.Session{
		ID:        uuid.New().String(),
		Username:  username,
		CreatedAt: now,
		UpdatedAt: now,
		Title:     userQuestion,
	}
	// 写入数据库
	res, err := dao.CreateSession(ctx, session)
	if err != nil || res <= 0 {
		return "", fmt.Errorf("create session failed: %v", err)
	}
	return session.ID, nil
}

// StreamMessage 流式会话 向当前会话输出
func StreamMessage(ctx context.Context, userQuestion, username, sessionID string, writer http.ResponseWriter, llmType llm.LLMType) error {
	// 流式输出需要http支持 flusher 它允许 HTTP 服务器将 部分响应数据立即发送到客户端
	// 而不必等待整个响应生成完成。这打破了传统 HTTP 请求-响应的"全量返回"模式，实现了 服务器推送 的效果。
	if _, ok := writer.(http.Flusher); !ok {
		return fmt.Errorf("writer does not support flushing")
	}
	manager := llm.GetGlobalAgentManager()
	agent, err := manager.GetOrCreateAgent(ctx, username, sessionID, llmType)
	if err != nil {
		return fmt.Errorf("get agent failed: %v", err)
	}
	// 定义一个回调函数支持SSE Server-Sent-Events
	// 	SSE的核心是“HTTP长连接的响应流”：
	// 客户端发起GET请求，服务器返回200 OK并设置text/event-stream响应头，且不关闭连接；
	// 服务器持续向客户端推送格式化文本（如data: 订单状态更新为已发货\n\n），每段数据以\n\n结尾，
	// 浏览器的EventSource会自动解析并触发onmessage事件；
	// 连接断开后，浏览器根据retry字段（默认3秒）自动重新发起请求，无需前端写重连代码。
	callback := func(msg string) error {
		log.Infof("[SSE] Sending Chunk: %s (length: %d)", msg, len(msg))
		// 写入响应体
		_, err = fmt.Fprintf(writer, "data: %s\n\n", msg)
		if err != nil {
			log.Errorf("[SSE] Write Response Failed: %v", err)
			return fmt.Errorf("write response failed: %v", err)
		}
		// 刷新响应体 确保立即发送到客户端
		writer.(http.Flusher).Flush()
		return nil
	}
	// 流式输出
	_, err = agent.StreamResponse(ctx, userQuestion, callback)
	if err != nil {
		log.Errorf("[SSE] Stream Response Failed: %v", err)
		return fmt.Errorf("stream response failed: %v", err)
	}
	// 发完了向客户端发送一个结束事件
	_, err = fmt.Fprintf(writer, "data: [DONE]\n\n")
	if err != nil {
		log.Errorf("[SSE] Write End Event Failed: %v", err)
		return fmt.Errorf("write end event failed: %v", err)
	}
	// 刷新响应体 确保立即发送到客户端
	writer.(http.Flusher).Flush()
	return nil
}

// ChatStreamSend 聊天 流式会话中返回对问题的解答
func ChatStreamSend(ctx context.Context, userQuestion, username, sessionID string, writer http.ResponseWriter, llmType llm.LLMType) error {
	return StreamMessage(ctx, userQuestion, username, sessionID, writer, llmType)
}

// CreateStreamSessionAndSendMessage 创建流式会话并发送消息
func CreateStreamSessionAndSendMessage(ctx context.Context, userQuestion, username string, writer http.ResponseWriter, llmType llm.LLMType) (string, error) {
	sessionID, err := CreateStreamSession(ctx, userQuestion, username)
	if err != nil {
		return "", fmt.Errorf("create stream session failed: %v", err)
	}
	err = StreamMessage(ctx, userQuestion, username, sessionID, writer, llmType)
	if err != nil {
		return "", fmt.Errorf("stream message failed: %v", err)
	}
	return sessionID, nil
}

// ChatHistory 获取当前会话的聊天历史记录
func ChatHistory(ctx context.Context, sessionID string) ([]model.MessageHistory, error) {
	messages, err := dao.GetMessagesBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get messages failed: %v", err)
	}
	msgHistory := make([]model.MessageHistory, 0, len(messages))
	for _, msg := range messages {
		msgHistory = append(msgHistory, model.MessageHistory{
			IsUser:    msg.IsUser,
			Username:  msg.Username,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
	}
	return msgHistory, nil
}
