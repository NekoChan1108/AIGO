package logic

import (
	"AIGO/internal/dao"
	"AIGO/internal/model"
	"AIGO/pkg/db"
	"context"
	"fmt"
)

const EmailVerificationCacheKey = "email:"

// Register 注册
/*
@description: 注册
@param ctx 上下文
@param username 用户名
@param email 邮箱
@param password 密码
@param verificationCode 验证码
@return error 错误
*/
func Register(ctx context.Context, username, email, password, verificationCode string) error {
	// 根据参数建立用户
	var user = &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	if verificationCode == "" {
		return fmt.Errorf("verificationCode can't be empty")
	}
	cmd := db.RedisDB.Get(ctx, EmailVerificationCacheKey+email)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("get verificationCode from redis failed: %v", err)
	}
	if cmd.Val() != verificationCode {
		return fmt.Errorf("verificationCode is invalid")
	}
	// 验证码校验成功后删除验证码
	go func() {
		cmd := db.RedisDB.Del(ctx, EmailVerificationCacheKey+email)
		if err := cmd.Err(); err != nil {
			fmt.Println("delete verificationCode from redis failed: ", err)
		}
	}()
	res, err := dao.Register(ctx, user)
	if err != nil || res <= 0 {
		return err
	}
	return nil
}

// Login 登录
/*
@description: 登录
@param ctx 上下文
@param username 用户名
@param email 邮箱
@param password 密码
@return *model.User 用户信息
@return error 错误
*/
func Login(ctx context.Context, username, email, password string) (*model.User, error) {
	var user = &model.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	return dao.Login(ctx, user)
}
