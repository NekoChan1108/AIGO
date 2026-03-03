package model

import "time"

type Session struct {
	ID        string    `bson:"_id" json:"id"`                // 会话ID
	Title     string    `bson:"title" json:"title"`           // 会话标题
	Username  string    `bson:"username" json:"username"`     // 用户名
	CreatedAt time.Time `bson:"created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"` // 更新时间
	DeletedAt time.Time `bson:"deleted_at" json:"deleted_at"` // 删除时间
}

// SessionInfo 返回给前端的会话信息
type SessionInfo struct {
	ID     string `json:"id"`     // 会话ID
	Title  string `json:"title"`  // 会话标题
}
