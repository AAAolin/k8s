package main

import (
	"github.com/gin-gonic/gin"
	"k8s/config"
	"k8s/controller"
	"k8s/db"
	"k8s/middle"
	"k8s/service"
)

func main() {
	// 初始化路由
	r := gin.Default()
	// 初始化clientSet
	service.K8s.Init()
	// 初始化数据库
	db.DB.Init()
	// 使用跨域、token认证中间件
	r.Use(middle.Cors())
	r.Use(middle.JWT())
	// 注册绑定路由
	controller.Router.InitApiRouter(r)
	// 启动
	r.Run(config.ListenAddr)
}
