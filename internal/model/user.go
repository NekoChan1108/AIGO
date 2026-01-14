package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户
type User struct {
	ID          uint           `gorm:"primaryKey" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	LastLoginAt *time.Time     `json:"last_login_at"` // 最后登录时间 采用指针接受空值避免初次创建报错
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	// gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password string `gorm:"not null;size:255" json:"-"` // json:"-" 不参与序列化不返回给前端 255为了存加密后的密码
}

// 指定表明
func (User) TableName() string {
	return "aigo_user"
}
