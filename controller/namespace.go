package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type namespace struct{}

var Namespace namespace

func (n *namespace) GetNamespaceList(c *gin.Context) {
	params := new(struct {
		FilterName string `json:"filterName"`
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
	data, err := service.Namespace.GetNamespaceList(params.FilterName, params.Page, params.Limit)
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

func (n *namespace) GetNamespaceDetail(c *gin.Context) {
	params := new(struct {
		NamespaceName string `json:"namespaceName"`
	})

	if err := c.Bind(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Namespace.GetNamespaceDetail(params.NamespaceName)
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

func (n *namespace) DeleteNamespace(c *gin.Context) {
	params := new(struct {
		NamespaceName string `json:"namespaceName"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("数据绑定失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Namespace.DeleteNamespace(params.NamespaceName)
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
