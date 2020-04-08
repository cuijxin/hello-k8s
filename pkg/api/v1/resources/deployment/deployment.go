package deployment

import (
	apps "k8s.io/api/apps/v1"
)

// CreateDeploymentRequest 定义了创建一个Deployment对象时所需的参数
type CreateDeploymentRequest struct {
	// Deployment 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// Deployment deployment对象定义.
	Deployment apps.Deployment `json:"deployment"`
}

// DeleteDeploymentRequest 定义了删除一个Deployment对象时所需参数.
type DeleteDeploymentRequest struct {
	// Name Deployment对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}
