package storageclass

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/storageclass"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 获取所有 StorageClass 对象列表.
// @Description 获取某一用户创建的所有Job对象
// @Tags resource
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/storageclass/list [get]
func GetStorageClassList(c *gin.Context) {
	log.Debug("调用获取 StorageClass 对象列表的函数.")

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	dsQuery := dataselect.NewDataSelectQuery(dataselect.NoPagination, dataselect.NoSort, dataselect.NoFilter, dataselect.NoMetrics)
	list, err := storageclass.GetStorageClassList(clientset, dsQuery)
	if err != nil {
		SendResponse(c, errno.ErrGetStorageClassList, err)
		return
	}

	SendResponse(c, errno.OK, list)
}