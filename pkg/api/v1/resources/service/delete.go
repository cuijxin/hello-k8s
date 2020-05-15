package service

import (
	"context"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"

	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 删除指定 Service 对象
// @Description 删除指定 Service 对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body service.DeleteServiceRequest true "删除 Service 对象时所需的参数."
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/service/delete [delete]
func Delete(c *gin.Context) {
	log.Info("调用删除 Service 对象的函数.")

	var r DeleteServiceRequest
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
	if err := clientset.CoreV1().Services(r.Namespace).Delete(context.TODO(), r.Name, options); err != nil {
		tool.SendResponse(c, errno.ErrDeleteService, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}
