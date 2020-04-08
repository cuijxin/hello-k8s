package mysql

// CreateOperatorRequest 安装MySQL Operator组件时所需参数.
type CreateOperatorRequest struct {
	// Namespace Mysql Operator 安装在哪个命名空间之下.
	Namespace string `json:"namespace"`
}

// DeleteOperatorRequest 删除MySQL Operator组件时所需参数.
type DeleteOperatorRequest struct {
	// Namespace Mysql Operator 组件安装在哪个命名空间之下.
	Namespace string `json:"namespace"`
}

// CreateClusterRequest 创建MySQL集群时所需参数.
type CreateClusterRequest struct {
	// Name MySQL集群名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// Template mysql cluster template.
	Template MySQLClusterTemplate `json:"template"`
}

// DeleteClusterRequest 定义了删除MySQL集群时所需参数
type DeleteClusterRequest struct {

	// MySQL 集群名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}

// MySQLClusterTemplate MySQL集群模版.
type MySQLClusterTemplate struct {
	// Members defines the number of MySQL instances in a cluster
	Members int32 `json:"members"`
}
