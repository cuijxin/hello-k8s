package common

// ConfigMapItem 定义了构建一个 ConfigMap 时每一个配置文件的配置项
type ConfigMapItem struct {
	// Key 配置项的 Key.
	Key string `json:"key"`

	// Value 配置项的 Value. 目前只支持base64格式编码的字符串.
	Value string `json:"value"`
}

// ConfigMapArg 用户自定义配置文件数据
type ConfigMapArg struct {
	// Name ConfigMap 对象名称
	// Name string `json:"name"`

	// Items ConfigMap 对象 配置项
	Items []ConfigMapItem `json:"items"`
}

// DataVolumeArg数据存储选项
type DataVolumeArg struct {
	// Name volumeClaimName
	// Name string `json:"name"`

	// StorageClassName 存储类名称
	StorageClassName *string `json:"storageClassName,omitempty"`

	// AccessModes 访问模式
	AccessModes []string `json:"accessModes,omitempty"`

	// Capacity 存储容量
	Capacity float64 `json:"capacity"`
}

// CustomRootPasswordArg 用户自定义root密码选项
type CustomRootPasswordArg struct {
	// SecretValue
	SecretValue string `json:"value"`
}
