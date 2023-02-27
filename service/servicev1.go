package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
)

type serviceV1 struct{}

var ServiceV1 serviceV1

type ServiceRes struct {
	Items []corev1.Service `json:"items"`
	Total int              `json:"total"`
}

type CreateService struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Labels     map[string]string `json:"labels"`
	Port       int32             `json:"port"`
	TargetPort int32             `json:"targetPort"`
	NodePort   int32             `json:"nodePort"`
	Type       string            `json:"type"`
}

// GetServiceList 列表
func (s *serviceV1) GetServiceList(namespace, filterName string, page, limit int) (*ServiceRes, error) {
	serviceList, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取service列表失败", err.Error())
		return nil, errors.New("获取service列表失败" + err.Error())
	}
	dataSelector := DataSelector{
		GenericDataList: s.toDataCell(serviceList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}
	service := dataSelector.Filter()
	total := len(service.GenericDataList)
	data := service.Sort().Paginate()

	return &ServiceRes{
		Items: s.toCoreV1(data.GenericDataList),
		Total: total,
	}, nil
}

// GetServiceDetail 获取Service详情
func (s *serviceV1) GetServiceDetail(serviceName, namespace string) (*corev1.Service, error) {
	data, err := K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取service详情失败", err.Error())
		return nil, errors.New("获取service详情失败" + err.Error())
	}
	return data, nil
}

// CreateService 创建Service
func (s *serviceV1) CreateService(cs *CreateService) (*corev1.Service, error) {
	fmt.Println(cs)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cs.Name,
			Namespace: cs.Namespace,
			Labels:    cs.Labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     cs.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: cs.TargetPort,
					},
				},
			},
			Selector: cs.Labels,
			Type:     corev1.ServiceType(cs.Type),
		},
	}
	if cs.NodePort != 0 && cs.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = cs.NodePort
	}
	data, err := K8s.ClientSet.CoreV1().Services(cs.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建service失败", err.Error())
		return nil, errors.New("创建service失败" + err.Error())
	}
	return data, nil
}

// DeleteService 删除Service
func (s *serviceV1) DeleteService(serviceName, namespace string) error {
	err := K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除service失败", err.Error())
		return errors.New("删除service失败" + err.Error())
	}
	return nil
}

// UpdateService 更新Service
func (s *serviceV1) UpdateService(content, namespace string) (*corev1.Service, error) {
	var service = &corev1.Service{}
	err := json.Unmarshal([]byte(content), service)
	if err != nil {
		logger.Error("反序列化失败", err.Error())
		return nil, errors.New("反序列化失败" + err.Error())
	}
	data, err := K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新service详情失败", err.Error())
		return nil, errors.New("更新service详情失败" + err.Error())
	}
	return data, nil
}

func (s *serviceV1) toDataCell(serviceList []corev1.Service) []DataCell {
	data := make([]DataCell, len(serviceList))
	for i, _ := range serviceList {
		data[i] = serviceCell(serviceList[i])
	}
	return data
}

func (s *serviceV1) toCoreV1(serviceList []DataCell) []corev1.Service {
	data := make([]corev1.Service, len(serviceList))
	for i, _ := range serviceList {
		data[i] = corev1.Service(serviceList[i].(serviceCell))
	}
	return data
}
