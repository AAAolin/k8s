package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"k8s/config"
)

type db struct{}

var DB db

var (
	GORM   *gorm.DB
	isInit bool
)

func (db *db) Init() {

	if isInit {
		return
	}

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%d)/%v?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPassword, config.DBIP, config.DBPort, config.DBName)
	var err error
	GORM, err = gorm.Open(config.DBType, dsn)
	if err != nil {
		logger.Error("数据库连接失败", err.Error())
	}

	// 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	// 设置了连接可复用的最大时间
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	// 设置了连接可复用的最大时间
	GORM.DB().SetConnMaxLifetime(config.MaxLifeTime)

	isInit = true

	logger.Info("数据库连接成功")

}
