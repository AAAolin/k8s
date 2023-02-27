package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type namespace struct{}

var Namespace namespace

type NamespaceRes struct {
	Items []corev1.Namespace `json:"items"`
	Total int                `json:"total"`
}

// GetNamespaceList 获取列表
func (n *namespace) GetNamespaceList(filterName string, page, limit int) (*NamespaceRes, error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败", err.Error())
		return nil, errors.New("获取namespace列表失败" + err.Error())
	}
	dataSelector := DataSelector{
		GenericDataList: n.toDataCell(namespaceList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}

	namespace := dataSelector.Filter()
	total := len(namespace.GenericDataList)
	data := namespace.Sort().Paginate()

	return &NamespaceRes{
		Items: n.toCorev1(data.GenericDataList),
		Total: total,
	}, nil
}

// GetNamespaceDetail 获取namespace详情
func (n *namespace) GetNamespaceDetail(namespaceName string) (*corev1.Namespace, error) {
	data, err := K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取namespace详情失败", err.Error())
		return nil, errors.New("获取namespace详情失败" + err.Error())
	}
	return data, nil
}

// DeleteNamespace 删除namespace
func (n *namespace) DeleteNamespace(namespaceName string) error {
	err := K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除namespace失败", err.Error())
		return errors.New("删除namespace失败" + err.Error())
	}
	return nil
}

func (n *namespace) toDataCell(namespaceList []corev1.Namespace) []DataCell {
	data := make([]DataCell, len(namespaceList))
	for i, _ := range namespaceList {
		data[i] = namespaceCell(namespaceList[i])
	}
	return data
}

func (n *namespace) toCorev1(namespaceList []DataCell) []corev1.Namespace {
	data := make([]corev1.Namespace, len(namespaceList))
	for i, _ := range namespaceList {
		data[i] = corev1.Namespace(namespaceList[i].(namespaceCell))
	}
	return data
}
