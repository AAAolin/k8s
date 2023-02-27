package service

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s/config"
)

type k8s struct {
	ClientSet *kubernetes.Clientset
}

var K8s k8s

func (k *k8s) Init() {
	// 将kubeconfig格式化为rest.config类型
	config, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		println("kubeconfig格式化为rest.config类型失败:", err)
	}
	// 通过config创建clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		println("通过config创建clientset失败:", err)
	}
	k.ClientSet = clientSet
}
