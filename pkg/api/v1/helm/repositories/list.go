package repositories

import (
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/cmd/helm/search"
	"k8s.io/klog"
)

const searchMaxScore = 25

// @Summary 获取Helm的Repo中的Chart列表
// @Description 获取Helm的Repo中的Chart列表
// @Tags helm
// @Accept json
// @Produce json
// @Param version query string false "chart version"
// @Param versions query string false "if true, all versions"
// @Param keyword query string false "search keyword"
// @Success 200 {object} tool.Response "{"code":200,"mmessage":"OK","data":{""}}"
// @Router /v1/helm/repositories/charts [get]
func ListRepoCharts(c *gin.Context) {
	klog.Info("调用获取Helm仓库中chart列表的方法")

	version := c.Query("version")
	versions := c.Query("versions")
	keyword := c.Query("keyword")

	klog.Infof("keyword:%s", keyword)

	// default stable
	if version == "" {
		version = ">0.0.0"
	}

	index, err := buildSearchIndex(version)
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	var res []*search.Result
	if keyword == "" {
		res = index.All()
	} else {
		res, err = index.Search(keyword, searchMaxScore, false)
	}

	search.SortScore(res)
	var versionsB bool
	if versions == "true" {
		versionsB = true
	}
	data, err := applyConstraint(version, versionsB, res)
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	chartList := make(RepoChartList, 0, len(data))
	for _, v := range data {
		chartList = append(chartList, RepoChartElement{
			Name:        v.Name,
			Version:     v.Chart.Version,
			AppVersion:  v.Chart.AppVersion,
			Description: v.Chart.Description,
		})
	}

	tool.SendResponse(c, errno.OK, chartList)
}
