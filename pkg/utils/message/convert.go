package message

import (
	"AIGO/internal/model"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

// ConvertToSchemaMessage 将对话上下文消息转换为AI模型接受的消息
func ConvertToSchemaMessages(messages []*model.Message) []*schema.Message {
	schemaMessages := make([]*schema.Message, 0, len(messages))
	for _, msg := range messages {
		role := schema.Assistant
		if msg.IsUser {
			role = schema.User
		}
		// 转换为AI模型接受的消息
		schemaMessages = append(schemaMessages, &schema.Message{
			Role:    role,
			Content: msg.Content,
			Extra: map[string]any{
				"id":         msg.ID,
				"created_at": msg.CreatedAt.Format("2006-01-02 15:04:05"),
				"username":   msg.Username,
				"session_id": msg.SessionId,
				"is_user":    msg.IsUser,
			},
		})
	}
	return schemaMessages
}

// ConvertToModelMessage 将模型接受的消息转换为对话上下文消息供数据库存储
func ConvertToModelMessage(username, sessionId string, message *schema.Message) *model.Message {
	var isUser bool
	if message.Role == schema.User {
		isUser = true
	}
	return &model.Message{
		ID:        uuid.New().String(),
		Username:  username,
		Content:   message.Content,
		SessionId: sessionId,
		IsUser:    isUser,
		CreatedAt: time.Now(),
	}
}
