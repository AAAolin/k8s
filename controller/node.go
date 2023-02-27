package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type node struct{}

var Node node

func (n *node) GetNodeList(c *gin.Context) {
	params := new(struct {
		FilterName string `json:"filterName"`
		Limit      int    `json:"limit"`
		Page       int    `json:"page"`
	})

	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Node.GetNodeList(params.FilterName, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取node列表成功",
		"data": data,
	})
}

func (n *node) GetNodeDetail(c *gin.Context) {
	params := new(struct {
		NodeName string `json:"nodeName"`
	})

	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Node.GetNodeDetail(params.NodeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取node详情成功",
		"data": data,
	})
}
