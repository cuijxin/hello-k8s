package deployment

import (
	"hello-k8s/pkg/api/v1/tool"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/dataselect"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"
	"hello-k8s/pkg/utils/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 查询某一 Deployment 对象控制的Pods列表
// @Description 查询某一 Deployment 对象控制的Pods列表
// @Tags resource
// @Accept json
// @Produce json
// @Param name path string true "Deployment 对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} tool.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/deployment/pods/{name}/{namespace} [get]
func GetDeploymentPods(c *gin.Context) {
	log.Info("调用获取 Deployment 对象的 Pods 列表函数")

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

	podList, err := deployment.GetDeploymentPods(clientset, nil, dsQuery, namespace, name)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetDeploymentPodsList, err)
		return
	}

	tool.SendResponse(c, errno.OK, podList)
}
