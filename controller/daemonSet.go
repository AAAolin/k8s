package controller

import (
	"github.com/gin-gonic/gin"
	"k8s/service"
	"net/http"
)

type daemonSet struct{}

var DaemonSet daemonSet

func (d *daemonSet) GetDaemonSetList(c *gin.Context) {
	params := new(struct {
		FilterName string `json:"filterName"`
		Namespace  string `json:"namespace"`
		Limit      int    `json:"limit"`
		Page       int    `json:"page"`
	})

	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	data, err := service.DaemonSet.GetDaemonSetList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet列表成功",
		"data": data,
	})
}

func (d *daemonSet) GetDaemonSetDetail(c *gin.Context) {
	params := new(struct {
		DaemonSetName string `json:"daemonSetName"`
		Namespace     string `json:"namespace"`
	})

	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	data, err := service.DaemonSet.GetDaemonSetDetail(params.DaemonSetName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet详情成功",
		"data": data,
	})
}

func (d *daemonSet) DeleteDaemonSet(c *gin.Context) {
	params := new(struct {
		DaemonSetName string `json:"daemonSetName"`
		Namespace     string `json:"namespace"`
	})

	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	err := service.DaemonSet.DeleteDaemonSet(params.DaemonSetName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除DaemonSet成功",
		"data": nil,
	})
}

func (d *daemonSet) UpdateDaemonSet(c *gin.Context) {
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
	})

	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	data, err := service.DaemonSet.UpdateDaemonSet(params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新DaemonSet列表成功",
		"data": data,
	})
}
