package dao

import (
	"AIGO/internal/model"
	"AIGO/pkg/db"
	"fmt"
	"time"
)

// Register 注册 TODO添加邮箱验证码逻辑以及密码加密存储
func Register(user *model.User) (int64, error) {
	// 先查看用户是否存在
	res := db.MysqlDB.Where("username=? or email=?", user.Username, user.Email).Find(&model.User{})
	if res.RowsAffected > 0 {
		return 0, fmt.Errorf("user existed")
	}
	// 不存在就创建 开启事务
	tx := db.MysqlDB.Begin()
	defer func() {
		// 回滚
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	// 创建用户
	res = tx.Create(user)
	if res.Error != nil || res.RowsAffected <= 0 {
		// 回滚
		tx.Rollback()
		return 0, fmt.Errorf("insert user failed: %v", res.Error)
	}
	affectRows := res.RowsAffected
	// commit
	res = tx.Commit()
	if res.Error != nil {
		// 回滚
		tx.Rollback()
		return 0, fmt.Errorf("insert user commit failed: %v", res.Error)
	}
	return affectRows, nil
}

// Login 登录 TODO 添加缓存逻辑以及jwt和密码加密逻辑
func Login(user *model.User) (*model.User, error) {
	dbUser := &model.User{}
	query := db.MysqlDB
	// 先查看用户是否存在
	if user.Email != "" {
		query = query.Where("email=?", user.Email)
	} else {
		query = query.Where("username=?", user.Username)
	}
	res := query.First(dbUser)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("user not existed")
	}
	// 密码校验
	if dbUser.Password != user.Password {
		return nil, fmt.Errorf("password error")
	}
	// 更新用户登录时间
	res = db.MysqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			res.Rollback()
		}
	}()
	res = res.Model(dbUser).Update("last_login_at", time.Now())
	if res.Error != nil || res.RowsAffected <= 0 {
		res.Rollback()
		return nil, fmt.Errorf("update user last_login_at failed: %v", res.Error)
	}
	res = res.Commit()
	if res.Error != nil {
		res.Rollback()
		return nil, fmt.Errorf("update user last_login_at commit failed: %v", res.Error)
	}
	// 返回用户信息
	return dbUser, nil
}
