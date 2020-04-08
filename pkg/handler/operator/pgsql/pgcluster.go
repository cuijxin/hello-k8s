package pgsql

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	. "hello-k8s/pkg/handler"

	acidv1 "github.com/cuijxin/postgres-operator-atom/pkg/apis/acid.zalan.do/v1"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
)

// @Summary CreateCluster deploy postgres cluster to the kubernetes cluster.
// @Description CreateCluster deploy postgres cluster to the kubernetes cluster.
// @Tags cluster
// @Accept json
// @Produce json
// @param data body pgsql.CreateClusterRequest true "Deploy pgsqloperator params"
// @Success 200 {object} pgsql.CreateClusterResponse "{"code":0,"message":"OK","data":{}}"
// @Router /cluster/pgsqlcluster [post]
func CreateCluster(c *gin.Context) {
	log.Info("Pgsql Cluster deploy function called.")

	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	customClientset, err := client.NewPostgresClientSet()
	if err != nil {
		SendResponse(c, errno.ErrCreatePgsqlClientSet, nil)
	}

	var r CreateClusterRequest
	if err := c.BindJSON(&r); err != nil {
		SendResponse(c, errno.ErrBind, err)
		return
	}

	CreateNamespace(r.Namespace, clientset)

	pgc := newPostgresCluster(r)
	// log.Debugf("postgres cluster data: %v", pgc)
	result, err := customClientset.AcidV1().Postgresqls(r.Namespace).Create(pgc)
	if err != nil {
		SendResponse(c, errno.ErrCreatePostgresCluster, err)
		return
	}
	// log.Debugf("postgres cluser list: %v", result)

	SendResponse(c, nil, result)
}

// @Summary DeleteCluster delete postgres cluster from the kubernetes cluster.
// @Description DeleteCluster delete postgres clsuter from the kubernetes cluster.
// @Tags cluster
// @Accept json
// @Produce json
// @param data body pgsql.DeleteClusterRequest true "Delete postgres cluster params"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":{}}"
// @Router /cluster/pgsqlcluster [delete]
func DeleteCluster(c *gin.Context) {
	log.Info("Pgsql Cluster delete function called.")

	customClientset, err := client.NewPostgresClientSet()
	if err != nil {
		SendResponse(c, errno.ErrCreatePgsqlClientSet, nil)
	}

	var r DeleteClusterRequest
	if err := c.BindJSON(&r); err != nil {
		SendResponse(c, errno.ErrBind, err)
		return
	}

	err = customClientset.AcidV1().Postgresqls(r.Namespace).Delete(r.ClusterName, &metav1.DeleteOptions{})
	if err != nil {
		SendResponse(c, errno.ErrDeletePostgresCluster, err)
		return
	}

	SendResponse(c, errno.OK, nil)
}

// @Summary GetCluster get a postgres cluster information.
// @Description GetCluster get a postgres cluster information.
// @Tags cluster
// @Accept json
// @Produce json
// @param name path string true "postgres cluster name".
// @param namespace path string true "namespace".
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /cluster/pgsqlcluster/detail/{name}/{namespace} [get]
func GetCluster(c *gin.Context) {
	log.Info("Pgsql Cluster get function called.")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	customClientset, err := client.NewPostgresClientSet()
	if err != nil {
		SendResponse(c, errno.ErrCreatePgsqlClientSet, nil)
	}

	result, err := customClientset.AcidV1().Postgresqls(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		SendResponse(c, errno.ErrGetPostgresCluster, err)
		return
	}

	SendResponse(c, errno.OK, result)
}

// @Summary GetClusterList get the list of the postgres cluster.
// @Description GetClusterList get the list of the postgres cluster.
// @Tags cluster
// @Accept json
// @Produce json
// @param namespace path string true "namespace"
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /cluster/pgsqlcluster/list/{namespace} [get]
func GetClusterList(c *gin.Context) {
	log.Info("Pgsql Cluster list function called.")

	namespace := c.Param("namespace")
	if namespace == "" {
		SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	customClientset, err := client.NewPostgresClientSet()
	if err != nil {
		SendResponse(c, errno.ErrCreatePgsqlClientSet, nil)
	}

	result, err := customClientset.AcidV1().Postgresqls(namespace).List(metav1.ListOptions{})
	if err != nil {
		SendResponse(c, errno.ErrGetPostgresClusterList, err)
		return
	}

	SendResponse(c, errno.OK, result)
}

// func createNamespace(namespace string, clientset kubernetes.Interface) {
// 	_, err := clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
// 	if errors.IsNotFound(err) {
// 		ns := &corev1.Namespace{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name: namespace,
// 			},
// 		}
// 		clientset.CoreV1().Namespaces().Create(ns)
// 	}
// }

func newPostgresCluster(r CreateClusterRequest) *acidv1.Postgresql {
	return &acidv1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.ClusterName,
			Namespace: r.Namespace,
		},
		Spec: acidv1.PostgresSpec{
			TeamID: r.TeamId,
			Volume: acidv1.Volume{
				Size:         r.Volume.Size,
				StorageClass: r.Volume.StorageClass,
			},
			NumberOfInstances: r.Replicas,
			Users:             initUsers(r.Users),
			Databases:         r.Databases,
			PostgresqlParam: acidv1.PostgresqlParam{
				PgVersion: r.PostgresqlParam.PgVersion,
			},
		},
	}
}

func initUsers(users map[string]UserFlags) map[string]acidv1.UserFlags {
	result := map[string]acidv1.UserFlags{}
	for key, value := range users {
		var tmp acidv1.UserFlags
		for _, str := range value {
			tmp = append(tmp, str)
		}
		result[key] = tmp
	}
	return result
}
