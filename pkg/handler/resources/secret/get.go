package secret

import (
	. "hello-k8s/pkg/api/v1"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/secret"
	"hello-k8s/pkg/utils/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary  查询某一 Secret 对象的详情
// @Description 查询某一 Secret 对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "Secret 对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} handler.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/secret/detail/{name}/{namespace} [get]
func GetSecret(c *gin.Context) {
	log.Debug("调用获取 Secret 对象详情的函数")

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

	secret, err := secret.GetSecretDetail(clientset, namespace, name)
	if err != nil {
		SendResponse(c, errno.ErrGetSecret, err)
		return
	}

	SendResponse(c, errno.OK, secret)
}
