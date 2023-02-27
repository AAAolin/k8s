package service

import (
	"fmt"
	"k8s/dao"
	"k8s/model"
)

type workflow struct{}

var Workflow workflow

type CreateWorkflow struct {
	Name        string                  `json:"name"`
	Namespace   string                  `json:"namespace"`
	Labels      map[string]string       `json:"labels"`
	Replicase   int32                   `json:"replicase"`
	Image       string                  `json:"image"`
	Port        int32                   `json:"port"`
	HealthCheck bool                    `json:"healthCheck"`
	HealthPath  string                  `json:"healthPath"`
	Cpu         string                  `json:"cpu"`
	Memory      string                  `json:"memory"`
	TargetPort  int32                   `json:"targetPort"`
	NodePort    int32                   `json:"nodePort"`
	Type        string                  `json:"type"`
	Host        map[string][]*HttpPaths `json:"host"`
}

func (w *workflow) GetList(name, namespace string, page, limit int) (*dao.WorkflowRes, error) {
	data, err := dao.Workflow.GetList(name, namespace, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) GetById(id int) (*model.Workflow, error) {
	data, err := dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) Add(data *CreateWorkflow) error {
	// 创建资源
	err := createResource(data)
	if err != nil {
		return err
	}

	// 如果不创建ingress资源，那么表中ingress字段应为空
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}

	// 表中新增数据
	workflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicase,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}
	_, err = dao.Workflow.Add(workflow)
	if err != nil {
		return err
	}
	return nil
}

func (w *workflow) DelById(id int) error {
	data, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}
	err = deleteResource(data)
	if err != nil {
		return err
	}
	err = dao.Workflow.DelById(id)
	if err != nil {
		return err
	}
	return nil
}

func createResource(data *CreateWorkflow) error {
	cd := &DeploymentMeta{
		Name:        data.Name,
		Namespace:   data.Namespace,
		Labels:      data.Labels,
		Replicase:   data.Replicase,
		Image:       data.Image,
		Port:        data.Port,
		HealthCheck: data.HealthCheck,
		HealthPath:  data.HealthPath,
		Cpu:         data.Cpu,
		Memory:      data.Memory,
	}
	err := Deployment.CreateDeployment(cd)
	if err != nil {
		return err
	}
	// 既然使用ingress用来暴露应用，就不需要service使用nodeport类型来暴露了，只要使用clusterip类型用来集群内部使用就够
	var serviceType string
	if data.Type == "Ingress" {
		serviceType = "ClusterIP"
	} else {
		serviceType = data.Type
	}
	cs := &CreateService{
		Name:       getServiceName(data.Name),
		Namespace:  data.Namespace,
		Labels:     data.Labels,
		Port:       data.Port,
		TargetPort: data.TargetPort,
		NodePort:   data.NodePort,
		Type:       serviceType,
	}
	_, err = ServiceV1.CreateService(cs)
	if err != nil {
		return err
	}

	if data.Type == "Ingress" {
		ci := &CreateIngress{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Labels:    data.Labels,
			Host:      data.Host,
		}
		_, err = Ingress.CreateIngress(ci)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteResource(data *model.Workflow) error {
	err := Deployment.DeleteDeployment(data.Name, data.Namespace)
	if err != nil {
		return err
	}
	fmt.Printf("删除%v成功 \n", data.Name)
	err = ServiceV1.DeleteService(getServiceName(data.Name), data.Namespace)
	if err != nil {
		return err
	}
	fmt.Printf("删除%v成功 \n", getServiceName(data.Name))
	if data.Type == "Ingress" {
		err = Ingress.DeleteIngress(getIngressName(data.Name), data.Namespace)
		if err != nil {
			return err
		}
		fmt.Printf("删除%v成功 \n", getIngressName(data.Name))
	}
	return nil
}

func getServiceName(name string) string {
	return name + "-svc"
}

func getIngressName(name string) string {
	return name + "-ing"
}
