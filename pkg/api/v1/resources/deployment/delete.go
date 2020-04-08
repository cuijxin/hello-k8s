package deployment

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 删除指定Deployment对象.
// @Description 删除指定Deployment对象.
// @Tags resource
// @Accept json
// @Produce json
// @param data body deployment.DeleteDeploymentRequest true "删除一个Deployment对象时所需参数."
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/deployment/delete [delete]
func Delete(c *gin.Context) {
	log.Info("调用删除 Deployment 对象的函数.")

	var r DeleteDeploymentRequest
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	deletePolicy := metav1.DeletePropagationBackground
	if err := clientset.AppsV1().Deployments(r.Namespace).Delete(r.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		tool.SendResponse(c, errno.ErrDeleteDeployment, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
