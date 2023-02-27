package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s/config"
)

type login struct{}

var Login login

func (l *login) Login(username, password string) (err error) {
	if username == config.Username && password == config.Password {
		return nil
	} else {
		logger.Error("用户名或密码错误")
		return errors.New("用户名或密码错误")
	}
	return nil
}
