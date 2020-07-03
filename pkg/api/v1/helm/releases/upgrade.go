package releases

import (
	"errors"
	"fmt"
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
	"io"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/klog"
)

// @Summary 升级Release版本
// @Description 升级Release版本
// @Tags helm
// @Accept json
// @Produce json
// @Param release path string true "release名称"
// @Param namespace path string true "命名空间"
// @Param chart query string true "chart名"
// @Param data body releases.ReleaseOptions true "release参数"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/namespaces/{namespace}/releases/{release} [put]
func UpgradeRelease(c *gin.Context) {
	klog.Info("调用更新Release的函数.")

	name := c.Param("release")
	namespace := c.Param("namespace")
	chart := c.Query("chart")
	var options ReleaseOptions
	err := c.BindJSON(&options)
	if err != nil {
		if err != io.EOF {
			klog.Infof("got error: %v", err)
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
		tool.SendResponse(c, errno.InternalServerError, errors.New(fmt.Sprintf("upgrade options can not be empty")))
		return
	}

	vals, err := mergeValues(options)
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
	client := action.NewUpgrade(actionConfig)
	client.Namespace = namespace

	cp, err := client.ChartPathOptions.LocateChart(chart, config.Settings)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	if req := chartRequested.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			klog.Infof("got error: %v", err)
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
	}

	_, err = client.Run(name, chartRequested, vals)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
