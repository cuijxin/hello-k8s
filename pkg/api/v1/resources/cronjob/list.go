package cronjob

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/common"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/cronjob"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 获取某一用户空间下的所有 CronJob 对象
// @Description 获取某一用户空间下的所有 CronJob 对象
// @Tags resource
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/cronjob/list/{namespace} [get]
func GetCronJobList(c *gin.Context) {
	log.Info("调用获取 CronJob 对象列表的函数")

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

	list, err := cronjob.GetCronJobList(clientset, namespaceQuery, dsQuery, nil)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetCronJobList, err)
		return
	}

	tool.SendResponse(c, errno.OK, list)
}
