package pod

import (
	. "hello-k8s/pkg/api/v1"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/common"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/pod"
	"hello-k8s/pkg/utils/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 获取某一命名空间下的所有 Pod 对象
// @Description 获取某一命名空间下的所有 Pod 对象
// @Tags resource
// @Param namespace path string true "命名空间"
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/pod/list/{namespace} [get]
func GetPodList(c *gin.Context) {
	log.Info("调用获取 Pod 对象列表的函数")

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

	list, err := pod.GetPodList(clientset, nil, namespaceQuery, dsQuery)
	if err != nil {
		SendResponse(c, errno.ErrGetPodList, err)
		return
	}

	SendResponse(c, errno.OK, list)
}
