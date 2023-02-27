package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type daemonSet struct{}

var DaemonSet daemonSet

type DaemonSetResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

// GetDaemonSetList 列表
func (d *daemonSet) GetDaemonSetList(filterName, namespace string, limit, page int) (*DaemonSetResp, error) {
	daemonSetList, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取daemonSet列表失败", err.Error())
		return nil, errors.New("获取daemonSet列表失败" + err.Error())
	}

	dataSelector := DataSelector{
		GenericDataList: d.toDataCell(daemonSetList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}

	daemonSet := dataSelector.Filter()
	total := len(daemonSet.GenericDataList)
	data := daemonSet.Sort().Paginate()
	items := d.toAppsV1DaemonSet(data.GenericDataList)

	return &DaemonSetResp{
		Items: items,
		Total: total,
	}, nil
}

// GetDaemonSetDetail 获取DaemonSet详情
func (d *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (*appsv1.DaemonSet, error) {
	daemonSet, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取daemonSet详情失败", err.Error())
		return nil, errors.New("获取daemonSet详情失败" + err.Error())
	}
	return daemonSet, nil
}

// DeleteDaemonSet 删除DaemonSet
func (d *daemonSet) DeleteDaemonSet(daemonSetName, namespace string) error {
	err := K8s.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除daemonSet失败", err.Error())
		return errors.New("删除daemonSet失败" + err.Error())
	}
	return nil
}

// UpdateDaemonSet 更新DaemonSet
func (d *daemonSet) UpdateDaemonSet(content, namespace string) (*appsv1.DaemonSet, error) {
	daemonSet := &appsv1.DaemonSet{}
	err := json.Unmarshal([]byte(content), daemonSet)
	if err != nil {
		logger.Error("反序列化失败", err.Error())
		return nil, errors.New("反序列化失败" + err.Error())
	}
	data, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新daemonSet失败", err.Error())
		return nil, errors.New("更新daemonSet失败" + err.Error())
	}
	return data, nil
}

func (d *daemonSet) toDataCell(daemonSetList []appsv1.DaemonSet) []DataCell {
	daemonSet := make([]DataCell, len(daemonSetList))
	for i, _ := range daemonSetList {
		daemonSet[i] = daemonSetCell(daemonSetList[i])
	}
	return daemonSet
}

func (d *daemonSet) toAppsV1DaemonSet(daemonSetList []DataCell) []appsv1.DaemonSet {
	daemonSet := make([]appsv1.DaemonSet, len(daemonSetList))
	for i, _ := range daemonSetList {
		daemonSet[i] = appsv1.DaemonSet(daemonSetList[i].(daemonSetCell))
	}
	return daemonSet
}
