package releases

import (
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/klog"
)

// @Summary 删除指定的Release
// @Description 删除指定的Release
// @Tags helm
// @Accept json
// @Produce json
// @Param release path string true "release名称"
// @Param namespace path string true "命名空间"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/namespaces/{namespace}/releases/{release} [delete]
func UnInstallRelease(c *gin.Context) {
	klog.Info("调用删除Release的函数.")

	name := c.Param("release")
	namespace := c.Param("namespace")
	klog.Infof("name is [%s] and namespace is [%s]", name, namespace)

	actionConfig, err := config.ActionConfigInit(namespace)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	client := action.NewUninstall(actionConfig)
	_, err = client.Run(name)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
