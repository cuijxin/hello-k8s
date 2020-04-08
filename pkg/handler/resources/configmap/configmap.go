package configmap

// ConfigMapItem 定义了构建一个 ConfigMap 时每一个配置文件的配置项.
type ConfigMapItem struct {
	// Key 配置项的 Key.
	Key string `json:"key"`

	// Value 配置项的 Value. 目前只支持base64格式编码的字符串.
	Value string `json:"value"`
}

// CreateConfigMapRequest 定义了创建一个ConfigMap对象时所需参数.
type CreateConfigMapRequest struct {
	// Name ConfigMap对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// ConfigMapItems 配置项数组.
	ConfigMapItems []ConfigMapItem `json:"item"`
}

// DeleteConfigMapRequest 定义了删除ConfigMap对象时所需参数
type DeleteConfigMapRequest struct {
	// ConfigMap 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}
