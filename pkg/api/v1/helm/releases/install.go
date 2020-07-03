package releases

import (
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
	"io"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"k8s.io/klog"
)

// @Summary 安装Release
// @Description 安装Release
// @Tags helm
// @Accept json
// @Produce json
// @Param namespace path string true "命名空间"
// @Param release path string true "release名称"
// @Param chart query string true "chart名"
// @Param data body releases.ReleaseOptions true "安装release时的参数"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/namespaces/{namespace}/releases/{release} [post]
func InstallRelease(c *gin.Context) {
	klog.Info("调用安装Release的函数.")

	name := c.Param("release")
	namespace := c.Param("namespace")
	chart := c.Query("chart")

	klog.Infof("namespace is: %s", namespace)

	var options ReleaseOptions
	err := c.BindJSON(&options)
	if err != nil && err != io.EOF {
		tool.SendResponse(c, errno.InternalServerError, err)
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
	client := action.NewInstall(actionConfig)
	client.ReleaseName = name
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

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          getter.All(config.Settings),
					RepositoryConfig: config.Settings.RepositoryConfig,
					RepositoryCache:  config.Settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					klog.Infof("got error: %v", err)
					tool.SendResponse(c, errno.InternalServerError, err)
					return
				}
			} else {
				klog.Infof("got error: %v", err)
				tool.SendResponse(c, errno.InternalServerError, err)
				return
			}
		}
	}

	_, err = client.Run(chartRequested, vals)
	if err != nil {
		klog.Infof("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
