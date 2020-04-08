package service

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/service"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 查询某一 Service 对象对应的Pods列表
// @Description 查询某一 Service 对象对应的Pods列表
// @Tags resource
// @Accept json
// @Produce json
// @Param name path string true "Service 对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} tool.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/service/pods/{name}/{namespace} [get]
func GetServicePods(c *gin.Context) {
	log.Info("调用获取 Service 对象对应的 Pods 列表函数.")

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

	dsQuery := dataselect.NewDataSelectQuery(dataselect.NoPagination, dataselect.NoSort, dataselect.NoFilter, dataselect.NoMetrics)

	podList, err := service.GetServicePods(clientset, nil, namespace, name, dsQuery)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetServicePodsList, err)
		return
	}

	tool.SendResponse(c, errno.OK, podList)
}
