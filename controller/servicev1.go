package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type servicev1 struct{}

var ServicV1 servicev1

func (s *servicev1) GetServiceList(c *gin.Context) {
	params := new(struct {
		FilterName string `json:"filterName"`
		Namespace  string `json:"namespace"`
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.ServiceV1.GetServiceList(params.Namespace, params.FilterName, params.Page, params.Limit)
	if err != nil {
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

func (s *servicev1) GetServiceDetail(c *gin.Context) {
	params := new(struct {
		Servicev1Name string `json:"servicev1name"`
		Namespace     string `json:"namespace"`
	})

	if err := c.Bind(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.ServiceV1.GetServiceDetail(params.Servicev1Name, params.Namespace)
	if err != nil {
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

func (s *servicev1) CreateService(c *gin.Context) {
	var ci = &service.CreateService{}
	if err := c.ShouldBindJSON(ci); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.ServiceV1.CreateService(ci)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建成功",
		"data": data,
	})
}

func (s *servicev1) DeleteService(c *gin.Context) {
	params := new(struct {
		Servicev1Name string `json:"servicev1Name"`
		Namespace     string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.ServiceV1.DeleteService(params.Servicev1Name, params.Namespace)
	if err != nil {
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

func (s *servicev1) UpdateService(c *gin.Context) {
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.ServiceV1.UpdateService(params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新成功",
		"data": data,
	})
}
