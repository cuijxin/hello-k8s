package releases

import (
	"errors"
	"fmt"
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/klog"
)

// @Summary 获取指定release的详细信息
// @Description 获取指定release的详细信息
// @Tags helm
// @Accept json
// @Produce json
// @Param namespace path string true "命名空间"
// @Param release path string true "release名"
// @Param info query string true "info:待获取信息的类型，可选值[all,hooks,manifest,notes,values]"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/namespaces/{namespace}/releases/{release} [get]
func ShowReleaseInfo(c *gin.Context) {
	klog.Info("调用获取指定Release详细信息的函数.")
	name := c.Param("release")
	namespace := c.Param("namespace")
	info := c.Query("info")
	infos := []string{"all", "hooks", "manifest", "notes", "values"}
	infoMap := map[string]bool{}
	for _, i := range infos {
		infoMap[i] = true
	}
	if _, ok := infoMap[info]; !ok {
		tool.SendResponse(c, errno.InternalServerError, errors.New(fmt.Sprintf("bad info %s, release info only support all/hooks/manifest/notes/values", info)))
		return
	}

	actionConfig, err := config.ActionConfigInit(namespace)
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	if info == "values" {
		client := action.NewGetValues(actionConfig)
		results, err := client.Run(name)
		if err != nil {
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
		tool.SendResponse(c, errno.OK, results)
		return
	}

	client := action.NewGet(actionConfig)
	results, err := client.Run(name)
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	if info == "all" {
		results.Chart = nil
		tool.SendResponse(c, errno.OK, results)
		return
	} else if info == "hooks" {
		if len(results.Hooks) < 1 {
			tool.SendResponse(c, errno.OK, []*release.Hook{})
			return
		}
		tool.SendResponse(c, errno.OK, results.Hooks)
		return
	} else if info == "notes" {
		tool.SendResponse(c, errno.OK, results.Info.Notes)
		return
	}
}
