package secret

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 删除指定Secret对象
// @Description 删除指定Secret对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body secret.DeleteSecretRequest true "删除参数"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/secret/delete [delete]
func Delete(c *gin.Context) {
	log.Debug("调用删除 Secret 对象函数")

	var r DeleteSecretRequest
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
	if err := clientset.CoreV1().Secrets(r.Namespace).Delete(r.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		tool.SendResponse(c, errno.ErrDeleteSecret, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
