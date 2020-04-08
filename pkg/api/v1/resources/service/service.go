package service

// DeleteServiceRequest 定义了删除一个 Service 对象时所需参数.
type DeleteServiceRequest struct {
	// Name Service 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间
	Namespace string `json:"namespace"`

	// ClusterID Kubernetes 集群ID.
	ClusterID string `json:"clusterId,omitempty"`
}
