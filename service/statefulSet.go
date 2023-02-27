package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type statefulSet struct{}

var StatefulSet statefulSet

type StatefulSetResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

// GetStatefulSetList 列表
func (d *statefulSet) GetStatefulSetList(filterName, namespace string, limit, page int) (*StatefulSetResp, error) {
	statefulSetList, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取statefulSet列表失败", err.Error())
		return nil, errors.New("获取statefulSet列表失败" + err.Error())
	}

	dataSelector := DataSelector{
		GenericDataList: d.toDataCell(statefulSetList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}

	statefulSet := dataSelector.Filter()
	total := len(statefulSet.GenericDataList)
	data := statefulSet.Sort().Paginate()
	items := d.toAppsV1StatefulSet(data.GenericDataList)

	return &StatefulSetResp{
		Items: items,
		Total: total,
	}, nil
}

// GetStatefulSetDetail 获取StatefulSet详情
func (d *statefulSet) GetStatefulSetDetail(statefulSetName, namespace string) (*appsv1.StatefulSet, error) {
	statefulSet, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取statefulSet详情失败", err.Error())
		return nil, errors.New("获取statefulSet详情失败" + err.Error())
	}
	return statefulSet, nil
}

// DeleteStatefulSet 删除StatefulSet
func (d *statefulSet) DeleteStatefulSet(statefulSetName, namespace string) error {
	err := K8s.ClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除statefulSet失败", err.Error())
		return errors.New("删除statefulSet失败" + err.Error())
	}
	return nil
}

// UpdateStatefulSet 更新StatefulSet
func (d *statefulSet) UpdateStatefulSet(content, namespace string) (*appsv1.StatefulSet, error) {
	statefulSet := &appsv1.StatefulSet{}
	err := json.Unmarshal([]byte(content), statefulSet)
	if err != nil {
		logger.Error("反序列化失败", err.Error())
		return nil, errors.New("反序列化失败" + err.Error())
	}
	data, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新statefulSet失败", err.Error())
		return nil, errors.New("更新statefulSet失败" + err.Error())
	}
	return data, nil
}

func (d *statefulSet) toDataCell(statefulSetList []appsv1.StatefulSet) []DataCell {
	statefulSet := make([]DataCell, len(statefulSetList))
	for i, _ := range statefulSetList {
		statefulSet[i] = statefulSetCell(statefulSetList[i])
	}
	return statefulSet
}

func (d *statefulSet) toAppsV1StatefulSet(statefulSetList []DataCell) []appsv1.StatefulSet {
	statefulSet := make([]appsv1.StatefulSet, len(statefulSetList))
	for i, _ := range statefulSetList {
		statefulSet[i] = appsv1.StatefulSet(statefulSetList[i].(statefulSetCell))
	}
	return statefulSet
}
