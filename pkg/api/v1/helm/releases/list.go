package releases

import (
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/klog"
)

// @Summary 获取releases列表
// @Description 获取releases列表
// @Tags helm
// @Accept json
// @Produce json
// @Param namespace path string true "命名空间"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/namespaces/{namespace}/releases [get]
func ListReleases(c *gin.Context) {
	klog.Info("调用获取Releases列表的函数.")

	namespace := c.Param("namespace")
	actionConfig, err := config.ActionConfigInit(namespace)
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	client := action.NewList(actionConfig)
	client.Deployed = true
	results, err := client.Run()
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	// Initialize the array so no results returns an empty array instead of null
	elements := make([]ReleaseElement, 0, len(results))
	for _, r := range results {
		elements = append(elements, constructReleaseElement(r, false))
	}

	tool.SendResponse(c, errno.OK, elements)
	return
}
