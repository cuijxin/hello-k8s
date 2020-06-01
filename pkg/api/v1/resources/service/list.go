package service

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/common"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/service"
	"hello-k8s/pkg/utils/errno"

	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 获取某一用户创建的所有 Service 对象
// @Description 获取某一用户创建的所有 Service 对象
// @Tags resource
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/resource/service/list/{namespace} [get]
func GetServiceList(c *gin.Context) {
	log.Info("调用获取 Service 对象列表的函数")

	namespace := c.Param("namespace")
	if namespace == "" {
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
	namespaceMap := make([]string, 0)
	namespaceMap = append(namespaceMap, namespace)
	namespaceQuery := common.NewNamespaceQuery(namespaceMap)

	list, err := service.GetServiceList(clientset, namespaceQuery, dsQuery)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetServiceList, err)
		return
	}

	tool.SendResponse(c, errno.OK, list)
}
