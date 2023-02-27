package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type ingress struct{}

var Ingress ingress

func (n *ingress) GetIngressList(c *gin.Context) {
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
	data, err := service.Ingress.GetIngressList(params.Namespace, params.FilterName, params.Page, params.Limit)
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

func (n *ingress) GetIngressDetail(c *gin.Context) {
	params := new(struct {
		IngressName string `json:"ingressName"`
		Namespace   string `json:"namespace"`
	})

	if err := c.Bind(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Ingress.GetIngressDetail(params.IngressName, params.Namespace)
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

func (n *ingress) CreateIngress(c *gin.Context) {
	var ci = &service.CreateIngress{}
	if err := c.ShouldBindJSON(ci); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Ingress.CreateIngress(ci)
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

func (n *ingress) DeleteIngress(c *gin.Context) {
	params := new(struct {
		IngressName string `json:"ingressName"`
		Namespace   string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Ingress.DeleteIngress(params.IngressName, params.Namespace)
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

func (n *ingress) UpdateIngress(c *gin.Context) {
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

	data, err := service.Ingress.UpdateIngress(params.Content, params.Namespace)
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
