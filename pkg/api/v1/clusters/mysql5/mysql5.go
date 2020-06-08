package mysql5

import (
	"hello-k8s/pkg/model/common"
)

// ClusterOptions MySQL5 集群参数.
type ClusterOptions struct {
	// Name MySQL 集群名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// ClusterID Kubernetes 集群ID.
	ClusterID string `json:"clusterId,omitempty"`

	// Members MySQL集群节点数.
	Members *int32 `json:"members"`

	// InitDBName 用户初始化数据库名称.
	InitDBName *string `json:"dbName,omitempty"`

	// IsExport 是否开放公开服务.
	IsExport bool `json:"isExport,omitempty"`

	// Config 用户自定义配置文件.
	Config *common.ConfigMapArg `json:"config,omitempty"`

	// RootPassword 用户自定义root密码.
	RootPassword *common.CustomRootPasswordArg `json:"rootpassword,omitempty"`

	// DataVolume 用户定义持久化存储参数.
	DataVolume *common.DataVolumeArg `json:"datavolume,omitempty"`

	// BackupVolume 用户自定义备份数据持久化存储参数.
	BackupVolume *common.DataVolumeArg `json:"backupvolume,omitempty"`
}

// ClusterInfo MySQL V5集群信息
type ClusterInfo struct {
	// Name MySQL集群名称
	Name string `json:"name"`

	// Namespace 命名空间
	Namespace string `json:"namespace"`

	// Kubernetes 集群ID
	ClusterID string `json:"clusterId"`

	// RootUserName
	RootUserName string `json:"root"`

	// RootPassword root密码
	RootPassword string `json:"password,omitempty"`

	// InitDBName 初始数据库
	InitDBName string `json:"dbname,omitempty"`

	// Host mysql服务名
	Host string `json:"host"`

	// Port mysql服务端口
	Port int32 `json:"port"`

	// Domain mysql服务访问域名
	Domain string `json:"domain"`

	// Status mysql集群运行状态
	Status common.ServiceStatus `json:"status,omitempty"`
}
