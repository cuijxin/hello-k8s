package atomservice

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// @Summary 更新Atom自定义服务的镜像.
// @Description 更新Atom自定义服务的镜像.
// @Tags atomapp
// @Accept json
// @Produce json
// @param data body atomservice.UpdateAtomServiceImage true "更新Atom自定义服务的镜像时所需参数."
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /atomapp/atomservice/updateimage [post]
func UpdateImage(c *gin.Context) {
	log.Info("调用更新 Atom 自定义服务的镜像的函数.")
	var r UpdateAtomServiceImage
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

	deploymentsClient := clientset.AppsV1().Deployments(r.Namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(r.Name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		result.Spec.Template.Spec.Containers[0].Image = r.Image
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		SendResponse(c, errno.ErrUpdateDeploymentImage, retryErr)
		return
	}

	SendResponse(c, errno.OK, nil)
}
