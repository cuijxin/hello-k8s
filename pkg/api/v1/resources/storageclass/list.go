package storageclass

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/storageclass"
	"hello-k8s/pkg/utils/errno"

	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

// @Summary 获取所有 StorageClass 对象列表.
// @Description 获取某一用户创建的所有Job对象
// @Tags resource
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/resource/storageclass/list [get]
func GetStorageClassList(c *gin.Context) {
	klog.Info("调用获取 StorageClass 对象列表的函数.")

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	dsQuery := dataselect.NewDataSelectQuery(dataselect.NoPagination, dataselect.NoSort, dataselect.NoFilter, dataselect.NoMetrics)
	list, err := storageclass.GetStorageClassList(clientset, dsQuery)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetStorageClassList, err)
		return
	}

	tool.SendResponse(c, errno.OK, list)
}
