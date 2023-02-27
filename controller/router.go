package controller

import (
	"github.com/gin-gonic/gin"
)

type router struct{}

var Router router

func (r *router) InitApiRouter(e *gin.Engine) {

	e.POST("/api/login", Login.Auth).
		POST("/api/user/update", Login.UpdatePassword).
		POST("/api/user/register", Login.Register).
		// Pod
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetContainerName).
		GET("/api/k8s/pod/log", Pod.GetContainerLog).
		GET("/api/k8s/pod/numns", Pod.GetPodNsNum).
		// Deployment
		GET("/api/k8s/deployments", Deployment.GetDeployment).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		GET("/api/k8s/deployment/numns", Deployment.GetDeploymentNsNum).
		// DaemonSet
		GET("/api/k8s/daemonSet", DaemonSet.GetDaemonSetList).
		GET("/api/k8s/daemonSet/detail", DaemonSet.GetDaemonSetDetail).
		DELETE("/api/k8s/daemonSet/del", DaemonSet.DeleteDaemonSet).
		PUT("/api/k8s/daemonSet/update", DaemonSet.UpdateDaemonSet).
		// StatefulSet
		GET("/api/k8s/stateful", StatefulSet.GetStatefulSetList).
		GET("/api/k8s/stateful/detail", StatefulSet.GetStatefulSetDetail).
		DELETE("/api/k8s/stateful/del", StatefulSet.DeleteStatefulSet).
		PUT("/api/k8s/stateful/update", StatefulSet.UpdateStatefulSet).
		// Node
		GET("/api/k8s/node", Node.GetNodeList).
		GET("/api/k8s/node/detail", Node.GetNodeDetail).
		// Namespace
		GET("/api/k8s/namespace", Namespace.GetNamespaceList).
		GET("/api/k8s/namespace/detail", Namespace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace/del", Namespace.DeleteNamespace).
		// Service
		GET("/api/k8s/service", ServicV1.GetServiceList).
		GET("/api/k8s/service/detail", ServicV1.GetServiceDetail).
		POST("/api/k8s/service/create", ServicV1.CreateService).
		DELETE("/api/k8s/service/del", ServicV1.DeleteService).
		PUT("/api/k8s/service/update", ServicV1.UpdateService).
		// Ingress
		GET("/api/k8s/ingress", Ingress.GetIngressList).
		GET("/api/k8s/ingress/detail", Ingress.GetIngressDetail).
		POST("/api/k8s/ingress/create", Ingress.CreateIngress).
		DELETE("/api/k8s/ingress/del", Ingress.DeleteIngress).
		PUT("/api/k8s/ingress/update", Ingress.UpdateIngress).
		//workflow操作
		GET("/api/k8s/workflows", Workflow.GetList).
		GET("/api/k8s/workflow/detail", Workflow.GetById).
		POST("/api/k8s/workflow/create", Workflow.Add).
		DELETE("/api/k8s/workflow/del", Workflow.DelById).
		// PV
		GET("/api/k8s/pv", PV.GetPVList).
		GET("/api/k8s/pv/detail", PV.GetPVDetail).
		DELETE("/api/k8s/pv/del", PV.DeletePV)

}
