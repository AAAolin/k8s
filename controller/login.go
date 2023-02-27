package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type login struct{}

var Login login

func (l *login) Auth(c *gin.Context) {
	params := new(struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "数据绑定失败",
			"data": nil,
		})
		return
	}
	// err := service.Login.Login(params.UserName, params.Password)
	err := service.User.Login(params.UserName, params.Password)
	if err != nil {
		logger.Error("用户名或密码错误", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "用户名或密码错误",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "登录成功",
		"data": nil,
	})
}

func (l *login) UpdatePassword(c *gin.Context) {
	params := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	err := service.User.UpdatePassword(params.Username, params.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "修改密码失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "修改密码成功",
		"data": nil,
	})
}

func (l *login) Register(c *gin.Context) {
	params := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	err := service.User.Register(params.Username, params.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建用户失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建用户成功",
		"data": nil,
	})
}
