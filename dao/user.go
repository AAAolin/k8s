package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s/db"
	"k8s/model"
)

type user struct{}

var User user

type result struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetByUsername 查询
func (u *user) GetByUsername(username string) (*model.User, error) {
	var info = &model.User{}
	tx := db.GORM.Where("username = ?", username).First(info)
	if tx.Error != nil {
		logger.Error("查询失败", tx.Error.Error())
		return nil, errors.New("查询失败" + tx.Error.Error())
	}
	return info, nil
}

func (u *user) Update(username, password string) error {
	tx := db.GORM.Model(model.User{}).Where("username = ?", username).Update(model.User{Password: password})
	if tx.Error != nil {
		logger.Error("修改失败", tx.Error.Error())
		return errors.New("修改失败" + tx.Error.Error())
	}
	return nil
}

func (u *user) Add(username, password string) error {
	tx := db.GORM.Create(model.User{Id: 3, Username: username, Password: password})
	if tx.Error != nil {
		logger.Error("创建失败", tx.Error.Error())
		return errors.New("创建失败" + tx.Error.Error())
	}
	return nil
}
