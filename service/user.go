package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s/dao"
)

type user struct{}

var User user

// Login 登录
func (u *user) Login(username, password string) error {
	data, err := dao.User.GetByUsername(username)
	if err != nil {
		return errors.New("调用查询接口，查询失败")
	}
	if password != data.Password {
		return errors.New("用户名或密码错误")
	}
	return nil
}

// UpdatePassword 忘记密码
func (u *user) UpdatePassword(username, password string) error {
	err := dao.User.Update(username, password)
	if err != nil {
		logger.Error("修改密码失败", err.Error())
		return errors.New("修改密码失败" + err.Error())
	}
	return nil
}

// Register 注册
func (u *user) Register(username, password string) error {
	err := dao.User.Add(username, password)
	if err != nil {
		logger.Error("注册失败", err.Error())
		return errors.New("注册失败" + err.Error())
	}
	return nil
}
