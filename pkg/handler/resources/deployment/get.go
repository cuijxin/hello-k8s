package deployment

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary  查询某一 Deployment 对象的详情
// @Description 查询某一 Deployment 对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "Deployment 对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} handler.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/deployment/detail/{name}/{namespace} [get]
func GetDeployment(c *gin.Context) {
	log.Info("调用创建 Deployment 对象的函数")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	result, err := deployment.GetDeploymentDetail(clientset, namespace, name)
	if err != nil {
		SendResponse(c, errno.ErrGetDeployment, err)
		return
	}

	SendResponse(c, errno.OK, result)
}
