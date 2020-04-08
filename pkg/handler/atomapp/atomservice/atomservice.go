package atomservice

// import (
// 	deploy "hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"
// )
//
// // CreateAtomServiceRequest 创建Atom自定义服务时所需参数.
// type CreateAtomServiceRequest struct {
//
// 	// ClusterID Kubernetes 集群ID.
// 	ClusterID string `json:"clusterId,omitempty"`
//
// 	// Spec Atom自定义服务定义参数.
// 	Spec *deploy.AppDeploymentSpec `json:"spec"`
// }

// CreateAtomServiceResponse 创建Atom自定义服务时返回参数.
type CreateAtomServiceResponse struct {

	// Namespace 命名空间
	Namespace string `json:"namespace"`

	// Deployment deployment对象名称.
	Deployment string `json:"deployment"`

	// Service service对象名称.
	Service string `json:"service"`
}

// ScaleAtomServiceRequest 弹性伸缩Atom自定义服务时所需参数.
type ScaleAtomServiceRequest struct {
	// ClusterID Kubernetes 集群ID.
	// ClusterID string `json:"clusterId,omitempty"`

	// Name Atom自定义服务名称.
	Name string `json:"name"`

	// Namespace Atom自定义服务所在命名空间.
	Namespace string `json:"namespace"`

	// Replicas 弹性伸缩参数.
	Replicas int32 `json:"replicas"`
}

// UpdateAtomServiceImage 更新Atom自定义服务的镜像.
type UpdateAtomServiceImage struct {
	// ClusterID Kubernetes 集群ID.
	// ClusterID string `json:"clusterId,omitempty"`

	// Name Atom自定义服务名称.
	Name string `json:"name"`

	// Namespace Atom自定义服务所在命名空间.
	Namespace string `json:"namespace"`

	// Image 镜像名称（repo/镜像名:tag）
	Image string `json:"image"`
}
