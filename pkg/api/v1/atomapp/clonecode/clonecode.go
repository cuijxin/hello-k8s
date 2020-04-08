package clonecode

// GithubAuth github auth
type GithubAuth struct {
	// Username github用户名.
	//Username string `json:"username"`

	// Token github用户token.
	Token string `json:"token"`
}

// CreateCloneCodeJobRequest 创建一个克隆代码的Job时所需参数.
type CreateCloneCodeJobRequest struct {
	// Name 克隆代码Job的名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// ClusterID Kubernetes 集群ID.
	// ClusterID string `json:"clusterId"`

	// GithubRepo github 仓库名称
	GithubRepoURL string `json:"githubrepourl"`

	// GithubRepoBranchOrTagName github 仓库分支名/tag名
	GithubRepoBranchOrTagName string `json:gitbranch,omitempty`

	// CodePersistentVolumeClaim 代码持久化卷声明.
	CodePersistentVolumeClaim string `json:"codevolume"`

	// GithubAuth github auth信息.
	GithubAuth GithubAuth `json:"githubauth"`
}
