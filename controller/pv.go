package controller

import (
	"github.com/gin-gonic/gin"
	"k8s/service"
	"net/http"
)

type pv struct{}

var PV pv

func (p *pv) GetPVList(c *gin.Context) {
	params := new(struct {
		FilterName string `json:"filterName"`
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
	})

	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定数据失败",
			"data": nil,
		})
		return
	}

	data, err := service.PV.GetPVList(params.FilterName, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取列表失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取列表成功",
		"data": data,
	})
}

func (p *pv) GetPVDetail(c *gin.Context) {
	params := new(struct {
		PVName string `json:"PVName"`
	})

	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	data, err := service.PV.GetPVDetail(params.PVName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取详情失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取详情成功",
		"data": data,
	})
}

func (p *pv) DeletePV(c *gin.Context) {
	params := new(struct {
		PVName string `json:"PVName"`
	})

	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定数据失败",
			"data": nil,
		})
		return
	}

	err := service.PV.DeletePV(params.PVName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除成功",
		"data": nil,
	})
}
