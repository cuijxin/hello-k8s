package job

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/job"
	"hello-k8s/pkg/utils/errno"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary  查询某一Job对象的详情
// @Description 查询某一Job对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "Job对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} handler.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/job/detail/{name}/{namespace} [get]
func GetJob(c *gin.Context) {
	log.Debug("调用查询某一 Job 对象的函数.")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	job, err := job.GetJobDetail(clientset, namespace, name)
	if err != nil {
		SendResponse(c, errno.ErrGetJob, err)
		return
	}

	SendResponse(c, errno.OK, job)
}
