package configmap

import (
	. "hello-k8s/pkg/api/v1"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 删除指定 ConfigMap 对象
// @Description 删除指定 ConfigMap 对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body configmap.DeleteConfigMapRequest true "删除参数"
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/configmap/delete [delete]
func Delete(c *gin.Context) {
	log.Info("调用删除 ConfigMap 对象的函数")

	var r DeleteConfigMapRequest
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
	if err := clientset.CoreV1().ConfigMaps(r.Namespace).Delete(r.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		SendResponse(c, errno.ErrDeleteConfigMap, err)
		return
	}

	SendResponse(c, errno.OK, nil)
}
