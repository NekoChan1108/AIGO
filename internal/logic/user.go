package logic

import (
	"AIGO/internal/dao"
	"AIGO/internal/model"
)

// Register 注册
/*
@description: 注册
@param username 用户名
@param email 邮箱
@param password 密码
@return error 错误
*/
func Register(username, email, password string) error {
	// 根据参数建立用户
	var user = &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	res, err := dao.Register(user)
	if err != nil || res <= 0 {
		return err
	}
	return nil
}

// Login 登录
/*
@description: 登录
@param username 用户名
@param email 邮箱
@param password 密码
@return *model.User 用户信息
@return error 错误
*/
func Login(username, email, password string) (*model.User, error) {
	var user = &model.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	return dao.Login(user)
}
