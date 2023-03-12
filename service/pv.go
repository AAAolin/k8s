package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type pv struct{}

var PV pv

type pvRes struct {
	Items []corev1.PersistentVolume `json:"items"`
	Total int                       `json:"total"`
}

func (p *pv) GetPVList(filterName string, page, limit int) (*pvRes, error) {
	pvList, err := K8s.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取pv列表失败", err.Error())
		return nil, errors.New("获取pv列表失败" + err.Error())
	}
	dataSelector := &DataSelector{
		GenericDataList: p.toDataCell(pvList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}
	filtered := dataSelector.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	return &pvRes{
		Items: p.toCorev1PV(data.GenericDataList),
		Total: total,
	}, nil
}

func (p *pv) GetPVDetail(pvName string) (*corev1.PersistentVolume, error) {
	data, err := K8s.ClientSet.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取pv详情失败", err.Error())
		return nil, errors.New("获取pv详情失败" + err.Error())
	}
	return data, nil
}

func (p *pv) DeletePV(pvName string) error {
	err := K8s.ClientSet.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除失败", err.Error())
		return errors.New("删除失败" + err.Error())
	}
	return nil
}

func (p *pv) toDataCell(data []corev1.PersistentVolume) []DataCell {
	pvList := make([]DataCell, len(data))
	for i, _ := range data {
		pvList[i] = pvCell(data[i])
	}
	return pvList
}

func (p *pv) toCorev1PV(data []DataCell) []corev1.PersistentVolume {
	pvList := make([]corev1.PersistentVolume, len(data))
	for i, _ := range data {
		pvList[i] = corev1.PersistentVolume(data[i].(pvCell))
	}
	return pvList
}
