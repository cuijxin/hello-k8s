package job

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 删除指定Job对象
// @Description 删除指定Job对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body job.DeleteJobRequest true "删除参数"
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/job/delete [delete]
func DeleteJob(c *gin.Context) {
	log.Info("调用删除 Job 对象的函数")

	var r DeleteJobRequest
	if err := c.BindJSON(&r); err != nil {
		SendResponse(c, errno.ErrBind, err)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	deletePolicy := metav1.DeletePropagationBackground
	if err := clientset.BatchV1().Jobs(r.Namespace).Delete(r.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		SendResponse(c, errno.ErrDeleteJob, err)
		return
	}

	SendResponse(c, errno.OK, nil)
}
