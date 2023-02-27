package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type pod struct{}

var Pod pod

func (p *pod) GetPods(c *gin.Context) {
	params := new(struct {
		FilterName string `json:"filter_name"`
		Namespace  string `json:"namespace"`
		Limit      int    `json:"limit"`
		Page       int    `json:"page"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}

func (p *pod) GetPodDetail(c *gin.Context) {
	params := new(struct {
		PodName   string `json:"podName"`
		Namespace string `json:"namespace"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod详情成功",
		"data": data,
	})
}

func (p *pod) DeletePod(c *gin.Context) {
	params := new(struct {
		PodName   string `json:"podName"`
		Namespace string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除Pod成功",
	})
}

func (p *pod) UpdatePod(c *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.UpdatePod(params.Namespace, params.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新Pod成功",
		"data": data,
	})
}

func (p *pod) GetContainerName(c *gin.Context) {
	params := new(struct {
		PodName   string `json:"podName"`
		Namespace string `json:"namespace"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetContainerName(params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取容器名称成功",
		"data": data,
	})
}

func (p *pod) GetContainerLog(c *gin.Context) {
	params := new(struct {
		ContainerName string `json:"containerName"`
		PodName       string `json:"podName"`
		Namespace     string `json:"namespace"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetContainerLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取容器日志成功",
		"data": data,
	})
}

func (p *pod) GetPodNsNum(c *gin.Context) {
	data, err := service.Pod.GetPodNsNum()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取每个命名空间的pod数量成功",
		"data": data,
	})
}
