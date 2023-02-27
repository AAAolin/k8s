package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type node struct{}

var Node node

type NodeResp struct {
	Items []corev1.Node `json:"items"`
	Total int           `json:"total"`
}

// GetNodeList 获取Node列表
func (n *node) GetNodeList(filterName string, limit, page int) (*NodeResp, error) {
	nodeList, err := K8s.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取node列表失败", err.Error())
		return nil, errors.New("获取node列表失败" + err.Error())
	}

	dataSelector := DataSelector{
		GenericDataList: n.toDataCell(nodeList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}

	node := dataSelector.Filter()
	total := len(node.GenericDataList)
	data := node.Sort().Paginate()
	nodeAppv1List := n.toCorev1Node(data.GenericDataList)
	return &NodeResp{
		Items: nodeAppv1List,
		Total: total,
	}, nil
}

// GetNodeDetail 获取Node详情
func (n *node) GetNodeDetail(nodeName string) (*corev1.Node, error) {
	node, err := K8s.ClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取node详情失败", err.Error())
		return nil, errors.New("获取node详情失败" + err.Error())
	}
	return node, nil

}

func (n node) toDataCell(nodeList []corev1.Node) []DataCell {
	node := make([]DataCell, len(nodeList))
	for i, _ := range nodeList {
		node[i] = nodeCell(nodeList[i])
	}
	return node
}

func (n node) toCorev1Node(nodeList []DataCell) []corev1.Node {
	node := make([]corev1.Node, len(nodeList))
	for i, _ := range nodeList {
		node[i] = corev1.Node(nodeList[i].(nodeCell))
	}
	return node
}