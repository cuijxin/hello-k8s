package service

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	. "hello-k8s/pkg/handler"

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
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/service/delete [delete]
func Delete(c *gin.Context) {
	log.Info("调用删除 Service 对象的函数.")

	var r DeleteServiceRequest
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
	if err := clientset.CoreV1().Services(r.Namespace).Delete(r.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		SendResponse(c, errno.ErrDeleteService, err)
		return
	}

	SendResponse(c, errno.OK, nil)
}
