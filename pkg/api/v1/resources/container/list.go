package container

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/container"
	"hello-k8s/pkg/utils/errno"

	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 获取某一 Pod 中的所有容器对象.
// @Description 获取某一 Pod 中的所有容器对象.
// @Tags resource
// @Param podId path string true "Pod ID"
// @Param namespace path string true "命名空间"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/pod/container/{podId}/{namespace} [get]
func GetPodContainers(c *gin.Context) {
	log.Debug("调用获取某一 Pod 对象中的所有Containers对象的函数.")

	podID := c.Param("podId")
	namespace := c.Param("namespace")
	if podID == "" || namespace == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	containers, err := container.GetPodContainers(clientset, namespace, podID)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetPodContainers, err)
		return
	}

	tool.SendResponse(c, errno.OK, containers)
}
