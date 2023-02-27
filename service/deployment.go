package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"strconv"
	"time"
)

type deployment struct{}

var Deployment deployment

type DeploymentResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

type DeploymentMeta struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels"`
	Replicase   int32             `json:"replicase"`
	Image       string            `json:"image"`
	Port        int32             `json:"port"`
	HealthCheck bool              `json:"healthCheck"`
	HealthPath  string            `json:"healthPath"`
	Cpu         string            `json:"cpu"`
	Memory      string            `json:"memory"`
}

type DeploymentNsNum struct {
	Namespace string `json:"namespace"`
	Total     int    `json:"total"`
}

// GetDeployment 获取Deployment列表
func (d *deployment) GetDeployment(filerName, namespace string, limit, page int) (*DeploymentResp, error) {
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取deployment列表失败", err.Error())
		return nil, errors.New("获取deployment列表失败" + err.Error())
	}

	dataSelector := &DataSelector{
		GenericDataList: d.toDataCell(deploymentList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filerName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}

	deployment := dataSelector.Filter()
	total := len(deployment.GenericDataList)
	data := deployment.Sort().Paginate()
	deploymentappv1List := d.toAppv1Deployment(data.GenericDataList)

	return &DeploymentResp{
		Items: deploymentappv1List,
		Total: total,
	}, nil
}

// GetDeploymentDetail 获取Deployment详情
func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (*appsv1.Deployment, error) {
	deployment, err := K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取deployment详情失败", err.Error())
		return nil, errors.New("获取deployment详情失败" + err.Error())
	}
	return deployment, nil
}

// ScaleDeployment 修改Deployment副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, num int) (int32, error) {
	replicase, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取deployment副本数失败", err.Error())
		return 0, errors.New("获取deployment副本数失败" + err.Error())
	}

	replicase.Spec.Replicas = int32(num)

	replicaseNew, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, replicase, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("修改deployment副本数失败", err.Error())
		return 0, errors.New("修改deployment副本数失败" + err.Error())
	}
	return replicaseNew.Spec.Replicas, nil
}

// CreateDeployment 创建Deployment
func (d *deployment) CreateDeployment(data *DeploymentMeta) error {
	fmt.Println(data)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicase,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: data.Port,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.Port,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.Port,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}

	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	_, err := K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建deployment失败", err.Error())
		return errors.New("创建deployment失败" + err.Error())
	}
	return nil
}

// DeleteDeployment 删除Deployment
func (d *deployment) DeleteDeployment(deploymentName, namespace string) error {
	err := K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除deployment失败", err.Error())
		return errors.New("删除deployment失败" + err.Error())
	}
	return nil
}

// RestartDeployment 重启Deployment
func (d *deployment) RestartDeployment(deploymentName, namespace string) (*appsv1.Deployment, error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	byteData, err := json.Marshal(patchData)
	if err != nil {
		logger.Error("序列化失败", err.Error())
		return nil, errors.New("序列化失败" + err.Error())
	}
	deployment, err := K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", byteData, metav1.PatchOptions{})
	if err != nil {
		logger.Error("重启deployment失败", err.Error())
		return nil, errors.New("重启deployment失败" + err.Error())
	}
	return deployment, nil
}

// UpdateDeployment 更新Deployment
func (d *deployment) UpdateDeployment(content, namespace string) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error("反序列化失败", err.Error())
		return nil, errors.New("反序列化失败" + err.Error())
	}
	deploymentNew, err := K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新deployment失败", err.Error())
		return nil, errors.New("更新deployment失败" + err.Error())
	}
	return deploymentNew, nil
}

// GetDeploymentNsNum 获取每个namespace的Deployment数量
func (d *deployment) GetDeploymentNsNum() ([]*DeploymentNsNum, error) {
	deploymentNsNumSlice := []*DeploymentNsNum{}
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace失败", err.Error())
		return nil, errors.New("获取namespace失败" + err.Error())
	}
	for i, _ := range namespaceList.Items {
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespaceList.Items[i].Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取deployment失败", err.Error())
			return nil, errors.New("获取deployment失败" + err.Error())
		}
		deploymentNsNum := &DeploymentNsNum{
			Namespace: namespaceList.Items[i].Name,
			Total:     len(deploymentList.Items),
		}
		deploymentNsNumSlice = append(deploymentNsNumSlice, deploymentNsNum)
	}
	return deploymentNsNumSlice, nil
}

func (d *deployment) toDataCell(deploymentList []appsv1.Deployment) []DataCell {
	deployment := make([]DataCell, len(deploymentList))
	for i, _ := range deploymentList {
		deployment[i] = deploymentCell(deploymentList[i])
	}
	return deployment
}

func (d *deployment) toAppv1Deployment(deploymentList []DataCell) []appsv1.Deployment {
	deployment := make([]appsv1.Deployment, len(deploymentList))
	for i, _ := range deploymentList {
		deployment[i] = appsv1.Deployment(deploymentList[i].(deploymentCell))
	}
	return deployment
}
