package storageclass

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/storageclass"
	"hello-k8s/pkg/utils/errno"

	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 查询某一 StorageClass 对象的详情.
// @Description 查询某一 StorageClass 对象的详情.
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "StorageClass 对象名称"
// @Success 200 {object} tool.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /v1/resource/storageclass/detail/{name} [get]
func GetStorageClass(c *gin.Context) {
	log.Debug("调用查询 StorageClass 对象的函数.")

	name := c.Param("name")
	if name == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	sc, err := storageclass.GetStorageClass(clientset, name)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetStorageClass, err)
		return
	}

	tool.SendResponse(c, errno.OK, sc)
}
