package secret

import (
	. "hello-k8s/pkg/api/v1"
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 创建 Secret 对象
// @Description 创建 Secret 对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body secret.CreateSecretRequest true "创建 Secret 对象时所需参数"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /resource/secret/create [post]
func Create(c *gin.Context) {
	log.Debug("调用创建 Secret 对象的函数")

	var r CreateSecretRequest
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

	CreateNamespace(r.Namespace, clientset)

	s := newSecret(r)
	result, err := clientset.CoreV1().Secrets(r.Namespace).Create(s)
	if err != nil {
		SendResponse(c, errno.ErrCreateSecret, err)
		return
	}

	SendResponse(c, errno.OK, result)
}

func newSecret(r CreateSecretRequest) *v1.Secret {
	secret := v1.Secret{
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
