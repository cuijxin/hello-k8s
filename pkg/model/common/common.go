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

type AppType string

const (
	AtomApp    AppType = "atomapp"
	MySQL      AppType = "mysql"
	MySQLV5    AppType = "mysqlv5"
	Redis      AppType = "redis"
	PostgreSQL AppType = "Postgresql"
	RabbitMQ   AppType = "RabbitMQ"
)

type PgSQLData struct {
	// TeamID
	TeamID *string `json:"teamID"`

	// DefaultRootSecretName
	DefaultRootSecretName *string `json:"secret"`

	IngressRouteName *string `json:"ingressroute"`

	// Domain
	ServiceDomain *string `json:"servicedomain"`
}

type RedisData struct {
	Mode *string `json:"mode"`

	CustomRootPasswordName *string `json:"secret"`

	DataVolume []string `json:"volumes"`

	// Domain
	ServiceDomain *string `json:"domain"`
}

type RabbitMQData struct {
	Replicas *int32 `json:"replicas"`

	StatefulsetName *string `json:"statefulset"`

	ServiceDomain *string `json:"domain"`
}

type MySQLAddOnData struct {
	// Members MySQL集群节点数
	Members *int32 `json:"members"`

	// Statefulset 对象名称
	StatefulsetName *string `json:"statefulset"`

	// Service 对象名称
	ServiceName *string `json:"service"`

	NodePortServiceName *string `json:"nodeportservice"`

	Port *int32 `json:"port"`

	// IngressRoute ingressroute对象
	IngressRouteName *string `json:"ingressroute,omitempty"`

	// Domain
	ServiceDomain *string `json:"servicedomain"`

	// Config 用户自定义configmap对象名称
	ConfigMapName *string `json:"config,omitempty"`

	// DataVolume 数据卷pvc对象名称
	DataVolumeName []string `json:"datavolume,omitempty"`

	// BackupVolume 备份卷pvc对象名称
	BackupVolumeName []string `json:"backupvolume,omitempty"`

	// RootPassword 密钥对象名称
	RootPasswordSecretName *string `json:"secret,omitempty"`
}

type AtomAppData struct {
	// Deployment 对象名称
	DeploymentName *string `json:"deployment"`

	// Service 对象名称
	ServiceName *string `json:"service"`

	// IngressRoute ingress router对象名称
	IngressRouterName *string `json:"ingressRouter,omitempty"`

	// Domain 服务域名
	Domain *string `json:"domain,omitempty"`

	// Hpa hpa对象名称
	HpaName *string `json:"hpa,omitempty"`

	// ConfigMapName 用户自定义configmap对象名称数组
	ConfigMapName []string `json:"config,omitempty"`

	// DataVolumeName 数据卷pvc对象名称
	DataVolumeName []string `json:"datavolume,omitempty"`

	// BuildImagePVCName 构建镜像时使用的PVC对象名称.
	BuildImagePVCName *string `json:"imagePVC,omitempty"`

	// CloneCodeJobName 克隆代码时的job名称.
	CloneCodeJobName *string `json:"cloneJob,omitempty"`

	// BuildImageJobName 构建镜像时的job名称.
	BuildImageJobName *string `json:"buildJob,omitempty"`
}

type AtomApplication struct {
	// Name 插件名称
	Name string `json:"name"`

	// Namespace 插件所属命名空间
	Namespace string `json:"namespace"`

	// ClusterID Kubernetes 集群ID
	ClusterID string `json:"clusterId"`

	// Type 插件类型
	Type AppType `json:"type"`

	// MySQLAddon
	MySQLAddon *MySQLAddOnData `json:"mysql,omitempty"`

	// PgSQLAddon
	PgSQLAddon *PgSQLData `json:"pgsql,omitempty"`

	RedisAddon *RedisData `json:"redis,omitempty"`

	RabbitmqAddon *RabbitMQData `json:"rabbitmq,omitempty"`

	// AtomApplication
	AtomApplication *AtomAppData `json:"atomapps,omitempty"`
}

type ServiceStatus string

const (
	ServiceRunning ServiceStatus = "running"

	ServiceWarnning ServiceStatus = "warnning"

	ServiceFailed ServiceStatus = "failed"
)
