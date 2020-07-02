package repositories

import (
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
	"sync"

	"helm.sh/helm/v3/pkg/repo"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

// @Summary 更新Helm的Repository信息
// @Description 更新Helm的Repository信息
// @Tags helm
// @Accept json
// @Produce json
// @Param data body config.HelmConfig true "更新Helm Repositories信息所需参数."
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/helm/repositories [put]
func UpdateRepositories(c *gin.Context) {
	klog.Info("调用更新Helm的Repositories的函数.")

	type errRepo struct {
		Name string
		Err  string
	}
	errRepoList := []errRepo{}

	var r config.HelmConfig
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
	}
	klog.Infof("request body is [%v]", r)

	var wg sync.WaitGroup
	for _, c := range r.HelmRepos {
		wg.Add(1)
		go func(c *repo.Entry) {
			defer wg.Done()
			err := updateChart(c)
			if err != nil {
				errRepoList = append(errRepoList, errRepo{
					Name: c.Name,
					Err:  err.Error(),
				})
			}
		}(c)
	}
	wg.Wait()

	if len(errRepoList) > 0 {
		tool.SendResponse(c, errno.InternalServerError, errRepoList)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
	return
}
