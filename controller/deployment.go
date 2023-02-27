package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s/service"
	"net/http"
)

type deployment struct{}

var Deployment deployment

func (d *deployment) GetDeployment(c *gin.Context) {
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
	data, err := service.Deployment.GetDeployment(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Deployment列表成功",
		"data": data,
	})
}

func (d *deployment) GetDeploymentDetail(c *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deploymentName"`
		Namespace      string `json:"namespace"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeploymentDetail(params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Deployment详情成功",
		"data": data,
	})
}

func (d *deployment) ScaleDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deploymentName"`
		Namespace      string `json:"namespace"`
		Num            int    `json:"num"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.ScaleDeployment(params.DeploymentName, params.Namespace, params.Num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "修改Deployment副本数成功",
		"data": data,
	})
}

func (d *deployment) CreateDeployment(c *gin.Context) {
	deploymentMeta := &service.DeploymentMeta{}
	if err := c.ShouldBindJSON(deploymentMeta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Deployment.CreateDeployment(deploymentMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "创建Deployment成功",
	})
}

func (d *deployment) DeleteDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deploymentName"`
		Namespace      string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err := service.Deployment.DeleteDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除Deployment成功",
	})
}

func (d *deployment) RestartDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deploymentName"`
		Namespace      string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("绑定数据失败", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.RestartDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "重启Deployment成功",
		"data": data,
	})
}

func (d *deployment) UpdateDeployment(c *gin.Context) {
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
	data, err := service.Deployment.UpdateDeployment(params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新Deployment成功",
		"data": data,
	})
}

func (d *deployment) GetDeploymentNsNum(c *gin.Context) {
	data, err := service.Deployment.GetDeploymentNsNum()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取每个命名空间的Deployment数量成功",
		"data": data,
	})
}
