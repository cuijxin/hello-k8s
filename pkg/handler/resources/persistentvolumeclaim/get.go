package persistentvolumeclaim

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"
	pvc "hello-k8s/pkg/kubernetes/kuberesource/resource/persistentvolumeclaim"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	. "hello-k8s/pkg/handler"
)

// @Summary  查询某一PersistentVolumeClaim对象的详情
// @Description 查询某一PersistentVolumeClaim对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "PersistentVolumeClaim对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} handler.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/persistentvolumeclaim/detail/{name}/{namespace} [get]
func GetPersistentVolumeClaim(c *gin.Context) {
	log.Info("调用获取 PersistentVolumeClaim 对象详情的函数")

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

	result, err := pvc.GetPersistentVolumeClaimDetail(clientset, namespace, name)
	if err != nil {
		SendResponse(c, errno.ErrGetPersistentVolumeClaim, err)
		return
	}

	SendResponse(c, errno.OK, result)
}
