package buildimage

// CreateBuildImageRequest 创建一个构建 Docker 镜像的 Pod 时所需参数
type CreateBuildImageRequest struct {
	// Name 构建 Docker 镜像的 Pod 的名称
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// Repository 镜像仓库路径
	Repository string `json:"repository`

	// ImageTag 镜像tag
	ImageTag string `json:"tag,omitempty"`

	// ContextPath 镜像构建上下文所在目录名
	ContextPath string `json:"context"`

	// PersistentVolumeClaimName 持久化卷声明名称.
	PersistentVolumeClaimName string `json:"pvcname"`
}
