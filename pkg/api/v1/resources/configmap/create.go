package configmap

import (
	"encoding/base64"
	. "hello-k8s/pkg/api/v1"
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 创建 ConfigMap 对象
// @Description 创建 ConfigMap 对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body configmap.CreateConfigMapRequest true "创建 ConfigMap 对象时所需参数"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /resource/configmap/create [post]
func Create(c *gin.Context) {
	log.Debug("调用创建 ConfigMap 对象的函数")

	var r CreateConfigMapRequest
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

	cm := newConfigMap(r)
	result, err := clientset.CoreV1().ConfigMaps(r.Namespace).Create(cm)
	if err != nil {
		SendResponse(c, errno.ErrCreateConfigMap, err)
		return
	}

	SendResponse(c, errno.OK, result)
}

func newConfigMap(r CreateConfigMapRequest) *v1.ConfigMap {
	log.Debug("init configmap object.")

	configmap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	if len(r.ConfigMapItems) > 0 {
		tmp := make(map[string]string)
		for _, item := range r.ConfigMapItems {
			d, _ := base64.StdEncoding.DecodeString(item.Value)
			str := String(d)
			tmp[item.Key] = str
		}
		configmap.Data = tmp
	}

	return &configmap
}