package configmap

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/configmap"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 查询某一 ConfigMap 对象的详情
// @Description 查询某一 ConfigMap 对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @Param name path string true "ConfigMap 对象名称"
// @Param namespace path string true "用户命名空间"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/resource/configmap/detail/{name}/{namespace} [get]
func GetConfigMap(c *gin.Context) {
	log.Debug("调用获取 ConfigMap 对象详情函数")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	detail, err := configmap.GetConfigMapDetail(clientset, namespace, name)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetConfigMapDetail, err)
		return
	}

	tool.SendResponse(c, errno.OK, detail)
}
