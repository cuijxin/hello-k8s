package secret

// SecretItem 定义了一个 Secret 对象的 secret 信息.
type SecretItem struct {
	// Key secret item key.
	Key string `json:"key"`

	// Value secret item value.
	Value string `json:"value"`
}

// CreateSecretRequest 定义了创建一个 Secret 对象时所需的参数.
type CreateSecretRequest struct {
	// Name Secret 对象的名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// SecretItems Secret 信息
	SecretItems []SecretItem `json:"item"`
}

// DeleteSecretRequest 定义了删除一个 Secret 对象时所需的参数.
type DeleteSecretRequest struct {
	// Name Secret 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}
