package postgres

type Volume struct {
	Size         string `json:"size"`
	StorageClass string `json:"storageClass,omitempty"`
}

type PostgresqlParam struct {
	PgVersion  *string           `json:"version"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type UserFlags []string

// CreateClusterOptions 定义了创建一个Postgresql集群时所需的参数.
type CreateClusterOptions struct {
	Name      string               `json:"name"`
	Namespace string               `json:"namespace"`
	ClusterID string               `json:"clusterId"`
	TeamID    *string              `json:"teamId"`
	Volume    *Volume              `json:"volume"`
	Replicas  *int32               `json:"replicas"`
	Users     map[string]UserFlags `json:"users"`
	Databases map[string]string    `json:"databases"`
	PgParams  *PostgresqlParam     `json:"pgparams"`
}

// ClusterInfo PgSQL集群信息
type ClusterInfo struct {
	// Name 用户指定的PgSQL集群名称
	Name string `json:"name"`

	// Namespace 命名空间
	Namespace string `json:"namespace"`

	// Kubernetes 集群ID
	ClusterID string `json:"clusterId"`

	// RootUser root用户名
	RootUser string `json:"root"`

	// RootPassword root用户密码
	RootPassword string `json:"password,omitempty"`

	// Host pgsql服务名
	Host string `json:"host"`

	// Port postgresql服务端口
	Port int32 `json:"port"`

	// Domain postgresql服务访问域名
	Domain string `json:"domain"`
}
