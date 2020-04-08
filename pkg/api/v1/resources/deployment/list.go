package deployment

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/common"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 获取某一用户创建的所有 Deployment 对象
// @Description 获取某一用户创建的所有 Deployment 对象
// @Tags resource
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/deployment/list/{namespace} [get]
func GetDeploymentList(c *gin.Context) {
	log.Info("调用获取 Deployment 对象列表的函数")

	namespace := c.Param("namespace")
	if namespace == "" {
		SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	dsQuery := dataselect.NewDataSelectQuery(dataselect.NoPagination, dataselect.NoSort, dataselect.NoFilter, dataselect.NoMetrics)
	namespaceMap := make([]string, 0)
	namespaceMap = append(namespaceMap, namespace)
	namespaceQuery := common.NewNamespaceQuery(namespaceMap)

	list, err := deployment.GetDeploymentList(clientset, namespaceQuery, dsQuery, nil)
	if err != nil {
		SendResponse(c, errno.ErrGetDeploymentList, err)
		return
	}

	SendResponse(c, errno.OK, list)
}
