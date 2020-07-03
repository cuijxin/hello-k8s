package releases

import (
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/klog"
)

// @Summary 获取Release历史版本信息
// @Description 获取Release历史版本信息
// @Tags helm
// @Accept json
// @Produce json
// @Param release path string true "release名称"
// @Param namespace path string true "命名空间"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/namespaces/{namespace}/releases/{release}/histories [get]
func ListReleaseHistories(c *gin.Context) {
	klog.Info("调用获取release历史版本信息的函数.")

	name := c.Param("release")
	namespace := c.Param("namespace")
	klog.Infof("name is [%s] and namespace is [%s]", name, namespace)

	actionConfig, err := config.ActionConfigInit(namespace)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	client := action.NewHistory(actionConfig)
	results, err := client.Run(name)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	if len(results) == 0 {
		tool.SendResponse(c, errno.OK, &ReleaseHistory{})
		return
	}

	tool.SendResponse(c, errno.OK, getReleaseHistory(results))
}
