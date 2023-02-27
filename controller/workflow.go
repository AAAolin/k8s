package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type workflow struct{}

var Workflow workflow

func (w *workflow) GetList(c *gin.Context) {
	params := new(struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Page      int    `json:"page"`
		Limit     int    `json:"limit"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Workflow.GetList(params.Name, params.Namespace, params.Page, params.Limit)
	if err != nil {
		logger.Error("获取失败", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功",
		"data": data,
	})
}
func (w *workflow) GetById(c *gin.Context) {
	params := new(struct {
		ID int `json:"id"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Workflow.GetById(params.ID)
	if err != nil {
		logger.Error("获取失败", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功",
		"data": data,
	})
}
func (w *workflow) Add(c *gin.Context) {
	var cf = &service.CreateWorkflow{}
	if err := c.ShouldBindJSON(cf); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Workflow.Add(cf)
	if err != nil {
		logger.Error("插入失败", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "插入成功",
		"data": nil,
	})
}
func (w *workflow) DelById(c *gin.Context) {
	params := new(struct {
		ID int `json:"id"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Workflow.DelById(params.ID)
	if err != nil {
		logger.Error("删除失败", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除成功",
		"data": nil,
	})
}
