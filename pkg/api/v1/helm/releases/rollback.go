package releases

import (
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
	"strconv"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/klog"
)

// @Summary 回滚Release版本
// @Description 回滚Release版本
// @Tags helm
// @Accept json
// @Produce json
// @Param release path string true "release名称"
// @Param namespace path string true "命名空间"
// @Param reversion path string true "版本"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/namespaces/{namespace}/releases/{release}/versions/{reversion} [put]
func RollBackRelease(c *gin.Context) {
	klog.Info("调用Release回滚操作的函数")

	name := c.Param("release")
	namespace := c.Param("namespace")
	reversionStr := c.Param("reversion")
	reversion, err := strconv.Atoi(reversionStr)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	actionConfig, err := config.ActionConfigInit(namespace)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	client := action.NewRollback(actionConfig)
	client.Version = reversion
	err = client.Run(name)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
