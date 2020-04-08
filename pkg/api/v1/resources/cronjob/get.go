package cronjob

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"hello-k8s/pkg/kubernetes/kuberesource/resource/cronjob"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary  查询某一 CronJob 对象的详情
// @Description 查询某一 CronJob 对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "CronJob 对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} tool.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/cronjob/detail/{name}/{namespace} [get]
func GetCronJob(c *gin.Context) {
	log.Debug("调用获取 CronJob 对象的函数")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	cronJob, err := cronjob.GetCronJobDetail(clientset, namespace, name)
	if err != nil {
		tool.SendResponse(c, errno.ErrGetCronJob, err)
		return
	}

	tool.SendResponse(c, errno.OK, cronJob)
}
