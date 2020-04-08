package atomservice

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 弹性伸缩Atom自定义服务的Pods数.
// @Description 弹性伸缩Atom自定义服务的Pods数.
// @Tags atomapp
// @Accept json
// @Produce json
// @param data body atomservice.ScaleAtomServiceRequest true "弹性伸缩Atom自定义服务时所需参数."
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /atomapp/atomservice/scale [post]
func Scale(c *gin.Context) {
	log.Info("调用弹性伸缩 Atom 自定义服务的函数.")
	var r ScaleAtomServiceRequest
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

	scale := autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: r.Namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: r.Replicas,
		},
	}

	result, err := clientset.AppsV1().Deployments(r.Namespace).UpdateScale(r.Name, &scale)
	if err != nil {
		SendResponse(c, errno.ErrScaleDeployment, err)
		return
	}

	SendResponse(c, errno.OK, result)
}