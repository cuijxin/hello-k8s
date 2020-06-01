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
