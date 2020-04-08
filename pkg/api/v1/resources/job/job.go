package job

import "hello-k8s/pkg/model"

// CreateJobRequest 定义了创建一个Job对象时所需参数.
type CreateJobRequest struct {
	// Name Job 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// Job 对象参数.
	Job model.JobArgs `json:"jobTemplate"`
}

// DeleteJobRequest 定义了删除Job对象时所需参数
type DeleteJobRequest struct {
	// Job 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}
