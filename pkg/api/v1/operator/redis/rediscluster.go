package redis

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"

	redisfailoverv1 "github.com/spotahome/redis-operator/api/redisfailover/v1"

	"hello-k8s/pkg/api/v1/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 创建RedisFailover集群.
// @Description 创建RedisFailover集群.
// @Tags cluster
// @Accept json
// @Produce json
// @param data body redis.CreateClusterRequest true "创建RedisFailover集群时所需参数"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{}}"
// @Router /cluster/rediscluster [post]
func CreateCluster(c *gin.Context) {
	log.Info("调用创建 RedisFailover 集群的函数.")

	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	customClientset, err := client.NewRedisClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateRedisClientSet, nil)
		return
	}

	var r CreateClusterRequest
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	tool.CreateNamespace(r.Namespace, clientset)

	redisObj := newRedisFailover(r)
	result, err := customClientset.DatabasesV1().RedisFailovers(r.Namespace).Create(redisObj)
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateRedisFailoverCluster, err)
		return
	}

	tool.SendResponse(c, nil, result)
}

// @Summary 从Kubernetes集群中删除已经部署的RedisFailover集群.
// @Description 从Kubernetes集群中删除已经部署的RedisFailover集群.
// @Tags cluster
// @Accept json
// @Produce json
// @param data body redis.DeleteClusterRequest true "删除RedisFailover集群时所需参数"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{}}"
// @Router /cluster/rediscluster [delete]
func DeleteCluster(c *gin.Context) {
	log.Info("调用删除RedisFailover集群的函数.")

	customClientset, err := client.NewRedisClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateRedisClientSet, nil)
		return
	}

	var r DeleteClusterRequest
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	err = customClientset.DatabasesV1().RedisFailovers(r.Namespace).Delete(r.Name, &metav1.DeleteOptions{})
	if err != nil {
		tool.SendResponse(c, errno.ErrDeleteRedisFailoverCluster, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}

// @Summary GetCluster get a redis cluster information.
// @Description GetCluster get a redis cluster information.
// @Tags cluster
// @Accept json
// @Produce json
// @param name path string true "redis cluster name".
// @param namespace path string true "namespace".
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /cluster/rediscluster/detail/{name}/{namespace} [get]
func GetCluster(c *gin.Context) {
	log.Info("Pgsql Cluster get function called.")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	customClientset, err := client.NewRedisClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateRedisClientSet, nil)
		return
	}

	result, err := customClientset.DatabasesV1().RedisFailovers(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		tool.SendResponse(c, errno.ErrGetRedisFailoverCluster, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

// @Summary GetClusterList get the list of the redis cluster.
// @Description GetClusterList get the list of the redis cluster.
// @Tags cluster
// @Accept json
// @Produce json
// @param namespace path string true "namespace"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /cluster/rediscluster/list/{namespace} [get]
func GetClusterList(c *gin.Context) {
	log.Info("Pgsql Cluster list function called.")

	namespace := c.Param("namespace")
	if namespace == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	customClientset, err := client.NewRedisClientSet()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateRedisClientSet, nil)
		return
	}

	result, err := customClientset.DatabasesV1().RedisFailovers(namespace).List(metav1.ListOptions{})
	if err != nil {
		tool.SendResponse(c, errno.ErrGetRedisFailoverClusterList, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

func newRedisFailover(r CreateClusterRequest) *redisfailoverv1.RedisFailover {
	redis := redisfailoverv1.RedisFailover{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	return &redis
}
