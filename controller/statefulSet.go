package controller

import (
	"github.com/gin-gonic/gin"
	"k8s/service"
	"net/http"
)

type statefulSet struct{}

var StatefulSet statefulSet

func (d *statefulSet) GetStatefulSetList(c *gin.Context) {
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
	data, err := service.StatefulSet.GetStatefulSetList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取statefulSet列表成功",
		"data": data,
	})
}

func (d *statefulSet) GetStatefulSetDetail(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"statefulSetName"`
		Namespace       string `json:"namespace"`
	})

	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	data, err := service.StatefulSet.GetStatefulSetDetail(params.StatefulSetName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取statefulSet详情成功",
		"data": data,
	})
}

func (d *statefulSet) DeleteStatefulSet(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"statefulSetName"`
		Namespace       string `json:"namespace"`
	})

	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	err := service.StatefulSet.DeleteStatefulSet(params.StatefulSetName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除statefulSet成功",
		"data": nil,
	})
}

func (d *statefulSet) UpdateStatefulSet(c *gin.Context) {
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
	data, err := service.StatefulSet.UpdateStatefulSet(params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新statefulSet列表成功",
		"data": data,
	})
}
