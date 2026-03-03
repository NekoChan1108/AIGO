package dao

import (
	"AIGO/internal/model"
	"AIGO/pkg/db"
	"AIGO/pkg/log"
	"AIGO/pkg/utils/encrypt"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// TODO 用户名密码邮箱的正则表达式校验 以及用户信息更新的功能


const UserCacheKey = "user:"

// Register 注册
func Register(ctx context.Context, user *model.User) (int64, error) {
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
	encryptPassword, err := encrypt.MD5Encrypt(user.Password)
	if err != nil {
		// 回滚
		tx.Rollback()
		return 0, err
	}
	user.Password = encryptPassword
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
	// commit成功之后添加缓存
	cacheUser, err := json.Marshal(user)
	if err != nil {
		return affectRows, fmt.Errorf("marshal user data failed: %v", err)
	}
	cmd := db.RedisDB.Set(ctx, UserCacheKey+user.Username, cacheUser, 0)
	if cmd.Err() != nil {
		return affectRows, fmt.Errorf("insert user cache failed: %v", cmd.Err())
	}
	return affectRows, nil
}

// Login 登录
func Login(ctx context.Context, user *model.User) (*model.User, error) {
	dbUser := &model.User{}
	query := db.MysqlDB
	// 先查看用户是否存在缓存 出错数据库兜底
	getCmd := db.RedisDB.Get(ctx, UserCacheKey+user.Username)
	if getCmd.Err() == nil {
		log.Info("hit cache")
		// 存在缓存中
		cacheUserBytes, err := getCmd.Bytes()
		if err == nil {
			// 缓存中存在
			if err := json.Unmarshal(cacheUserBytes, dbUser); err == nil {
				return dbUser, nil
			}
		}
	}
	log.Info("search db")
	// 不在缓存中 或者 缓存查询出错
	if user.Email != "" {
		query = query.Where("email=?", user.Email)
	} else {
		query = query.Where("username=?", user.Username)
	}
	res := query.First(dbUser)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("user not existed")
	}
	encryptPassword, err := encrypt.MD5Encrypt(user.Password)
	if err != nil {
		return nil, err
	}
	// 密码校验
	if dbUser.Password != encryptPassword {
		return nil, fmt.Errorf("password error")
	}
	// 更新用户登录时间
	now := time.Now()
	res = db.MysqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			res.Rollback()
		}
	}()
	res = res.Model(dbUser).Update("last_login_at", now)
	if res.Error != nil || res.RowsAffected <= 0 {
		res.Rollback()
		return nil, fmt.Errorf("update user last_login_at failed: %v", res.Error)
	}
	res = res.Commit()
	if res.Error != nil {
		res.Rollback()
		return nil, fmt.Errorf("update user last_login_at commit failed: %v", res.Error)
	}
	dbUser.LastLoginAt = &now
	go func() {
		// 更新缓存
		setCache(ctx, dbUser)
	}()
	// 返回用户信息
	return dbUser, nil
}

// setCache 设置缓存
func setCache(ctx context.Context, user *model.User) error {
	log.Info("set user cache")
	cacheUser, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal user data failed: %v", err)
	}
	cmd := db.RedisDB.Set(ctx, UserCacheKey+user.Username, cacheUser, 0)
	if cmd.Err() != nil {
		return fmt.Errorf("insert user cache failed: %v", cmd.Err())
	}
	return nil
}
