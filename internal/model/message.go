package model

import (
	"time"
)

// Message 人机交互的消息
type Message struct {
	ID        string    `bson:"_id" json:"id"`                // 消息ID
	IsUser    bool      `bson:"is_user" json:"is_user"`       // 标识这条消息是用户发送的还是模型发送的
	Content   string    `bson:"content" json:"content"`       // 消息内容
	Username  string    `bson:"username" json:"username"`     // 用户名
	SessionId string    `bson:"session_id" json:"session_id"` // 会话ID 标识这条消息属于哪个会话
	CreatedAt time.Time `bson:"created_at" json:"created_at"` // 创建时间
}

// MessageHistory 消息历史记录
type MessageHistory struct {
	IsUser    bool      `json:"is_user"`    // 标识这条消息是用户发送的还是模型发送的
	Username  string    `json:"username"`   // 用户名
	Content   string    `json:"content"`    // 消息历史内容
	CreatedAt time.Time `json:"created_at"` // 创建时间
}
