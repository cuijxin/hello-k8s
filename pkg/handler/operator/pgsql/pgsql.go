package pgsql

type CreateOpeartorRequest struct {
	ConfigMapName           string `json:"ConfigMapName"`
	DeploymentName          string `json:"DeploymentName"`
	Namespace               string `json:"Namespace"`
	ServiceAccountName      string `json:"ServiceAccountName`
	ClusterRoleName         string `json:"ClusterRoleName"`
	ClusterRoleBindingName  string `json:"ClusterRoleBindingName"`
	PostgresqlOperatorImage string `json:"PostgresqlOperatorImage"`
}

type CreateOperatorResponse struct {
	ConfigMapName          string `json:"ConfigMapName"`
	ServiceAccountName     string `json:"ServiceAccountName`
	ClusterRoleName        string `json:"ClusterRoleName"`
	ClusterRoleBindingName string `json:"ClusterRoleBindingName"`
	DeploymentName         string `json:"DeploymentName"`
	Namespace              string `json:"Namespace"`
}

type Volume struct {
	Size         string `json:"Size"`
	StorageClass string `json:"StorageClass"`
	SubPath      string `json:"SubPath,omitempty"`
}

type PostgresqlParam struct {
	PgVersion  string            `json:"version"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// UserFlags defines flags (such as superuser, nologin) that could be assigned to individual users
type UserFlags []string

type CreateClusterRequest struct {
	ClusterName     string               `json:"ClusterName"`
	Namespace       string               `json:"Namespace"`
	TeamId          string               `json:"TeamId"`
	Volume          Volume               `json:"Volume"`
	Replicas        int32                `json:"Replicas"`
	Users           map[string]UserFlags `json:"Users"`
	Databases       map[string]string    `json:"Databases"`
	PostgresqlParam PostgresqlParam      `json:"Postgresql"`
}

type CreateClusterResponse struct {
	ClusterName     string               `json:"ClusterName"`
	Namespace       string               `json:"Namespace"`
	TeamId          string               `json:"TeamId"`
	Volume          Volume               `json:"Volume"`
	Replicas        int32                `json:"Replicas"`
	Users           map[string]UserFlags `json:"Users"`
	Databases       map[string]string    `json:"Databases"`
	PostgresqlParam PostgresqlParam      `json:"Postgresql"`
}

type DeleteClusterRequest struct {
	ClusterName string `json:"ClusterName"`
	Namespace   string `json:"Namespace"`
}
