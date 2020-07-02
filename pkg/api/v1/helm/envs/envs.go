package envs

import (
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

// @Summary 获取指定Helm的Envs信息
// @Description 获取指定Helm的Envs信息
// @Tags helm
// @Accept json
// @Produce json
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/envs [get]
func GetHelmEnvs(c *gin.Context) {
	klog.Info("调用获取Helm的Envs信息的函数.")

	tool.SendResponse(c, errno.OK, config.Settings.EnvVars())

	return
}
