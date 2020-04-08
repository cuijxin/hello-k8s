package mysql

import (
	"hello-k8s/pkg/api/v1/tool"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
)

const (
	defaultAgentServiceAccoutName = "mysql-agent"
	defaultRoleName               = "mysql-agent"
	defaultAgentClusterRoleName   = "mysql-agent"
)

// @Summary 创建MySQL集群
// @Description 创建MySQL集群
// @Tags cluster
// @Accept json
// @Produce json
// @param data body mysql.CreateClusterRequest true "创建 MySQL 集群所需参数."
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /cluster/mysqlcluster [post]
func CreateCluster(c *gin.Context) {
	log.Debug("调用创建 MySQL 集群函数")

	var r CreateClusterRequest
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	customClientset, err := client.NewMySQLClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateMySQLClientSet, nil)
	}

	tool.CreateNamespace(r.Namespace, clientset)

	err = CheckMySQLClusterRBAC(
		r.Namespace,
		defaultAgentServiceAccoutName,
		defaultRoleName,
		defaultAgentClusterRoleName,
		clientset)
	if err != nil {
		tool.SendResponse(c, errno.ErrMySQLRBACCheck, err)
		return
	}

	mysqlCluster := newMySQLCluster(r)
	result, err := customClientset.MySQLV1alpha1().Clusters(r.Namespace).Create(mysqlCluster)
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateMySQLCluster, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

// @Summary 删除 Kubernetes 集群中的指定的 MySQL 集群.
// @Description 删除 Kubernetes 集群中的指定的 MySQL 集群.
// @Tags cluster
// @Accept json
// @Produce json
// @param data body mysql.DeleteClusterRequest true "删除 MySQL 集群时所需的参数."
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{}}"
// @Router /cluster/mysqlcluster [delete]
func DeleteCluster(c *gin.Context) {
	log.Debug("调用删除 MySQL 集群的函数.")

	customClientset, err := client.NewMySQLClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateMySQLClientSet, nil)
	}

	var r DeleteClusterRequest
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	deletePolicy := metav1.DeletePropagationBackground
	if err := customClientset.MySQLV1alpha1().Clusters(r.Namespace).Delete(r.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Error("Delete MySQL Cluster object failed.", err)
		tool.SendResponse(c, errno.ErrDeleteMySQLCluster, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}

// @Summary 查询指定的 MySQL 集群信息.
// @Description 查询指定的 MySQL 集群信息.
// @Tags cluster
// @Accept json
// @Produce json
// @param name path string true "postgres cluster name".
// @param namespace path string true "namespace".
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /cluster/mysqlcluster/detail/{name}/{namespace} [get]
func GetCluster(c *gin.Context) {
	log.Debug("调用查询指定的 MySQL 集群信息的函数.")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	customClientset, err := client.NewMySQLClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateMySQLClientSet, nil)
	}

	result, err := customClientset.MySQLV1alpha1().Clusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		tool.SendResponse(c, errno.ErrGetMySQLCluster, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

// @Summary 获取某一命名空间下的 MySQL 集群列表.
// @Description 获取某一命名空间下的 MySQL 集群列表.
// @Tags cluster
// @Accept json
// @Produce json
// @param namespace path string true "namespace"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /cluster/mysqlcluster/list/{namespace} [get]
func GetClusterList(c *gin.Context) {
	log.Debug("获取某一命名空间下的 MySQL 集群列表.")
	namespace := c.Param("namespace")
	if namespace == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	customClientset, err := client.NewMySQLClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateMySQLClientSet, nil)
	}
	result, err := customClientset.MySQLV1alpha1().Clusters(namespace).List(metav1.ListOptions{})
	if err != nil {
		tool.SendResponse(c, errno.ErrGetMySQLClusterList, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

func newMySQLCluster(r CreateClusterRequest) *v1alpha1.Cluster {
	mysqlCluster := v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	if reflect.ValueOf(r.Template).FieldByName("Members").IsValid() {
		mysqlCluster.Spec.Members = r.Template.Members
	}

	return &mysqlCluster
}
