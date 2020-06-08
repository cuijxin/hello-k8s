package secret

import (
	"context"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// @Summary 创建 Secret 对象
// @Description 创建 Secret 对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body secret.CreateSecretRequest true "创建 Secret 对象时所需参数"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/resource/secret/create [post]
func Create(c *gin.Context) {
	klog.Info("调用创建 Secret 对象的函数")

	var r CreateSecretRequest
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

	tool.CreateNamespace(r.Namespace, clientset)

	s := newSecret(r)
	result, err := clientset.CoreV1().Secrets(r.Namespace).Create(context.TODO(), s, metav1.CreateOptions{})
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateSecret, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

func newSecret(r CreateSecretRequest) *corev1.Secret {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	if len(r.SecretItems) > 0 {
		tmp := make(map[string]string)
		for _, item := range r.SecretItems {
			tmp[item.Key] = item.Value
		}
		secret.StringData = tmp
	}

	secret.Type = "Opaque"

	return &secret
}
