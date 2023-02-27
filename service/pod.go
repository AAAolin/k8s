package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s/config"
)

var Pod pod

type pod struct{}

// PodsResp 定义pod列表的返回内容，Items是corev1.pod类型的pod列表，Total为pod元素总数
//先过滤，再拿total，再做分页
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

type PodNsNum struct {
	Namespace string `json:"namespace"`
	Total     int    `json:"total"`
}

// GetPods 获取pod列表
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	//通过clintset获取pod完整列表
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取pod列表失败", err)
		return nil, errors.New("获取pod列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(podList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{Name: filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	filtered := selectableData.Filter()
	//再拿Total
	total := len(filtered.GenericDataList)
	//再排序和分页
	data := filtered.Sort().Paginate()
	//再将DataCell切片数据转成原生pod切片
	pods := p.fromCells(data.GenericDataList)
	//返回
	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}

// GetPodDetail 获取Pod详情
func (p *pod) GetPodDetail(podName, namespace string) (*corev1.Pod, error) {
	pod, err := K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取pod详情失败", err.Error())
		return nil, errors.New("获取pod详情失败" + err.Error())
	}
	return pod, nil
}

// DeletePod 删除Pod
func (p *pod) DeletePod(podName, namespace string) error {
	err := K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除pod失败", err.Error())
		return errors.New("删除pod失败" + err.Error())
	}
	return nil
}

// UpdatePod 更新Pod
func (p *pod) UpdatePod(namespace string, content string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := json.Unmarshal([]byte(content), pod)
	if err != nil {
		logger.Error("反序列化失败", err.Error())
		return nil, errors.New("反序列化失败" + err.Error())
	}
	podNew, err := K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("pod更新失败", err.Error())
		return nil, errors.New("pod更新失败" + err.Error())
	}
	return podNew, nil
}

// GetContainerName 获取Pod中的容器名
func (p *pod) GetContainerName(podName, namespace string) ([]string, error) {
	pod, _ := p.GetPodDetail(podName, namespace)
	containerName := []string{}

	for _, value := range pod.Spec.Containers {
		containerName = append(containerName, value.Name)
	}
	return containerName, nil

}

// GetContainerLog 获取容器日志
func (p *pod) GetContainerLog(containerName, podName, namespace string) (string, error) {
	limit := int64(config.LogsLimit)
	option := &corev1.PodLogOptions{
		Container:  containerName,
		LimitBytes: &limit,
	}
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)

	logs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error("获取log失败", err.Error())
		return "", errors.New("获取log失败" + err.Error())
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		logger.Error("复制PodLog失败", err)
		return "", errors.New("复制PodLog失败" + err.Error())
	}
	return buf.String(), nil
}

// GetPodNsNum 获取每个namespace的pod数量
func (p *pod) GetPodNsNum() ([]*PodNsNum, error) {
	podNsNumSlice := []*PodNsNum{}
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败", err)
		return nil, errors.New("获取namespace列表失败" + err.Error())
	}
	for _, namespace := range namespaceList.Items {
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取pod列表失败", err)
			return nil, errors.New("获取pod列表失败" + err.Error())
		}

		podNsNum := &PodNsNum{
			Namespace: namespace.Name,
			Total:     len(podList.Items),
		}
		podNsNumSlice = append(podNsNumSlice, podNsNum)
	}

	return podNsNumSlice, nil
}

//把podCell转成corev1 pod
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

//把corev1 pod转成DataCell
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}
