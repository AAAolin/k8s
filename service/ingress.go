package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type ingress struct{}

var Ingress ingress

type IngressRes struct {
	Items []nwv1.Ingress `json:"items"`
	Total int            `json:"total"`
}

type CreateIngress struct {
	Name      string                  `json:"name"`
	Namespace string                  `json:"namespace"`
	Labels    map[string]string       `json:"labels"`
	Host      map[string][]*HttpPaths `json:"host"`
}

type HttpPaths struct {
	Path        string        `json:"path,omitempty"`
	PathType    nwv1.PathType `json:"pathType,omitempty"`
	ServiceName string        `json:"serviceName,omitempty"`
	ServicePort int32         `json:"servicePort,omitempty"`
}

// GetIngressList 列表
func (i *ingress) GetIngressList(namespace, filterName string, page, limit int) (*IngressRes, error) {
	IngressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取Ingress列表失败", err.Error())
		return nil, errors.New("获取Ingress列表失败" + err.Error())
	}
	dataSelector := DataSelector{
		GenericDataList: i.toDataCell(IngressList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}
	Ingress := dataSelector.Filter()
	total := len(Ingress.GenericDataList)
	data := Ingress.Sort().Paginate()

	return &IngressRes{
		Items: i.tonwv1(data.GenericDataList),
		Total: total,
	}, nil
}

// GetIngressDetail 获取Ingress详情
func (i *ingress) GetIngressDetail(IngressName, namespace string) (*nwv1.Ingress, error) {
	data, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), IngressName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Ingress详情失败", err.Error())
		return nil, errors.New("获取Ingress详情失败" + err.Error())
	}
	return data, nil
}

// CreateIngress 创建Ingress
func (i *ingress) CreateIngress(ci *CreateIngress) (*nwv1.Ingress, error) {
	var ingressRules []nwv1.IngressRule
	var httpIngressPATHs []nwv1.HTTPIngressPath
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ci.Name,
			Namespace: ci.Namespace,
			Labels:    ci.Labels,
		},
	}
	for key, value := range ci.Host {
		ir := nwv1.IngressRule{
			Host: key,
			//这里现将nwv1.HTTPIngressRuleValue类型中的Paths置为空，后面组装好数据再赋值
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			//将每个hip对象组装成数组
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}
		//给Paths赋值，前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		//将每个ir对象组装成数组，这个ir对象就是IngressRule，每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules对象加入到ingress的规则中
	ingress.Spec.Rules = ingressRules
	data, err := K8s.ClientSet.NetworkingV1().Ingresses(ci.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建Ingress失败", err.Error())
		return nil, errors.New("创建Ingress失败" + err.Error())
	}
	return data, nil
}

// DeleteIngress 删除Ingress
func (i *ingress) DeleteIngress(IngressName, namespace string) error {
	err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), IngressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Ingress失败", err.Error())
		return errors.New("删除Ingress失败" + err.Error())
	}
	return nil
}

// UpdateIngress 更新Ingress
func (i *ingress) UpdateIngress(content, namespace string) (*nwv1.Ingress, error) {
	var Ingress = &nwv1.Ingress{}
	err := json.Unmarshal([]byte(content), Ingress)
	if err != nil {
		logger.Error("反序列化失败", err.Error())
		return nil, errors.New("反序列化失败" + err.Error())
	}
	data, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), Ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Ingress详情失败", err.Error())
		return nil, errors.New("更新Ingress详情失败" + err.Error())
	}
	return data, nil
}

func (i *ingress) toDataCell(IngressList []nwv1.Ingress) []DataCell {
	data := make([]DataCell, len(IngressList))
	for i, _ := range IngressList {
		data[i] = ingressCell(IngressList[i])
	}
	return data
}

func (i *ingress) tonwv1(IngressList []DataCell) []nwv1.Ingress {
	data := make([]nwv1.Ingress, len(IngressList))
	for i, _ := range IngressList {
		data[i] = nwv1.Ingress(IngressList[i].(ingressCell))
	}
	return data
}
