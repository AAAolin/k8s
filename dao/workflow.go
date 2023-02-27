package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s/db"
	"k8s/model"
)

type workflow struct{}

var Workflow workflow

type WorkflowRes struct {
	Items []*model.Workflow `json:"items"`
	Total int               `json:"total"`
}

func (w *workflow) GetList(name, namespace string, page, limit int) (*WorkflowRes, error) {
	var workflowList []*model.Workflow
	var total int
	startSet := limit * (page - 1)

	tx := db.GORM.
		Model(model.Workflow{}).
		Where("name like ? and namespace = ?", "%"+name+"%", namespace).
		Count(&total).
		Find(&workflowList)
	if startSet > total {
		startSet = 0
	}
	//tx := db.GORM.Where("namespace = ?", namespace).Find(&workflowList)
	//if tx.Error != nil {
	//	logger.Error("获取workflow列表失败", tx.Error)
	//	return nil, errors.New("获取workflow列表失败" + tx.Error.Error())
	//}
	//数据库查询，Limit方法用于限制条数，Offset方法设置起始位置
	tx = db.GORM.
		Model(model.Workflow{}).
		Where("name like ? and namespace = ?", "%"+name+"%", namespace).
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)
	if tx.Error != nil {
		logger.Error("查询失败", tx.Error.Error())
		return nil, errors.New("查询失败" + tx.Error.Error())
	}
	workflowRes := &WorkflowRes{
		Items: workflowList,
		Total: total,
	}
	return workflowRes, nil
}

func (w *workflow) GetById(id int) (*model.Workflow, error) {
	var workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).First(workflow)
	if tx.Error != nil {
		logger.Error("查询单条失败", tx.Error.Error())
		return nil, errors.New("查询单条失败" + tx.Error.Error())
	}
	return workflow, nil
}

func (w *workflow) Add(workflow *model.Workflow) (*model.Workflow, error) {
	tx := db.GORM.Create(workflow)
	if tx.Error != nil {
		logger.Error("插入失败", tx.Error.Error())
		return nil, errors.New("插入失败" + tx.Error.Error())
	}
	return workflow, nil
}

func (w *workflow) DelById(id int) error {
	//var workflow *model.Workflow  使用改方式初始化报错
	var workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).Delete(&workflow)
	if tx.Error != nil {
		logger.Error("删除失败", tx.Error.Error())
		return errors.New("删除失败" + tx.Error.Error())
	}
	return nil
}
