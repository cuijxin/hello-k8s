package persistentvolumeclaim

// CreatePersistentVolumeClaimRequest 定义了创建一个PersistentVolumeClaim对象时所需参数.
type CreatePersistentVolumeClaimRequest struct {
	// Name PersistentVolumeClaim对象名称
	Name string `json:"name"`
	// Namespace 命名空间
	Namespace string `json:"namespace"`

	// StoraegClassName 存储类名称.
	StorageClassName *string `json:"storageClassName"`

	// StorageCapacity 申请存储容量.
	StorageCapacity float64 `json:"storageCapacity"`

	// AccessModes 存储的访问模式.
	AccessModes []string `json:"AccessModes"`
}

// DeletePersistentVolumeClaimRequest 定义了删除PersistentVolumeClaim对象时所需参数
type DeletePersistentVolumeClaimRequest struct {
	// Kubernetes 集群ID.
	// ClusterID string `json:"clusterId"`

	// Secret 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}
