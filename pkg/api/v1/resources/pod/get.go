package pod

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/pod"
	"hello-k8s/pkg/utils/errno"

	"hello-k8s/pkg/api/v1/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary  查询某一 Pod 对象的详情
// @Description 查询某一 Pod 对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "Pod 对象名称"
// @Param namespace path string true "命名空间"
// @Success 200 {object} tool.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/pod/detail/{name}/{namespace} [get]
func GetPod(c *gin.Context) {
	log.Debug("调用获取 Pod 对象详情的函数.")

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

	pod, err := pod.GetPodDetail(clientset, nil, namespace, name)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetPodDetail, err)
		return
	}

	tool.SendResponse(c, errno.OK, pod)
}
