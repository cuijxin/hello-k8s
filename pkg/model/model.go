package model

type PostgresqlOperatorInfo struct {
	ConfigMapName           string `json:"ConfigMapName"`
	DeploymentName          string `json:"DeploymentName"`
	Namespace               string `json:"Namespace"`
	ServiceAccountName      string `json:"ServiceAccountName`
	ClusterRoleName         string `json:"ClusterRoleName"`
	ClusterRoleBindingName  string `json:"ClusterRoleBindingName"`
	PostgresqlOperatorImage string `json:"cuijx/postgres-operator:94a1a62"`
}
