package persistentvolumeclaim

import (
	"context"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 删除指定的PersistentVolumeClaim对象
// @Description 删除指定的PersistentVolumeClaim对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body persistentvolumeclaim.DeletePersistentVolumeClaimRequest true "删除参数"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/persistentvolumeclaim/delete [delete]
func Delete(c *gin.Context) {
	log.Info("调用删除PersistentVolumeClaim对象的函数")

	var r DeletePersistentVolumeClaimRequest
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

	deletePropagation := metav1.DeletePropagationBackground
	options := metav1.DeleteOptions{
		PropagationPolicy: &deletePropagation,
	}
	if err := clientset.CoreV1().PersistentVolumeClaims(r.Namespace).Delete(context.TODO(), r.Name, options); err != nil {
		tool.SendResponse(c, errno.ErrDeletePersistentVolumeClaim, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
