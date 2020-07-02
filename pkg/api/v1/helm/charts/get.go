package charts

import (
	"errors"
	"fmt"
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/klog"
)

// @Summary 获取指定chart的详细信息
// @Description 获取指定chart的详细信息
// @Tags helm
// @Accept json
// @Produce json
// @Param chart query string true "指定chart的name"
// @Param info query string true "指定想要获取的信息：readme/values/chart"
// @Param version query string false "指定想要获取信息的chart的版本"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/charts [get]
func GetChartInfo(c *gin.Context) {
	klog.Info("调用获取指定chart详细信息的方法")

	name := c.Query("chart")
	info := c.Query("info") // readme, values, chart
	version := c.Query("version")

	client := action.NewShow(action.ShowAll)
	client.Version = version
	if info == string(action.ShowChart) {
		client.OutputFormat = action.ShowChart
	} else if info == string(action.ShowReadme) {
		client.OutputFormat = action.ShowReadme
	} else if info == string(action.ShowValues) {
		client.OutputFormat = action.ShowValues
	} else {
		tool.SendResponse(c, errno.InternalServerError, errors.New(fmt.Sprintf("bad info %s, chart info only support readme/values/chart", info)))
		return
	}

	cp, err := client.ChartPathOptions.LocateChart(name, config.Settings)
	if err != nil {
		klog.Fatalf("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	chrt, err := loader.Load(cp)
	if err != nil {
		klog.Fatalf("got error: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	if client.OutputFormat == action.ShowChart {
		tool.SendResponse(c, errno.OK, chrt.Metadata)
		return
	}
	if client.OutputFormat == action.ShowValues {
		values := make([]*file, 0, len(chrt.Raw))
		for _, v := range chrt.Raw {
			values = append(values, &file{
				Name: v.Name,
				Data: string(v.Data),
			})
		}
		tool.SendResponse(c, errno.OK, values)
		return
	}
	if client.OutputFormat == action.ShowReadme {
		tool.SendResponse(c, errno.OK, string(findReadme(chrt.Files).Data))
		return
	}
}
