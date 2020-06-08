package persistentvolumeclaim

import (
	"context"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"strconv"

	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// @Summary 创建PersistentVolumeClaim对象
// @Description 创建PersistentVolumeClaim对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body persistentvolumeclaim.CreatePersistentVolumeClaimRequest true "创建PersistentVolumeClaim对象所需参数."
// @Success 200 {object} tool.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /v1/resource/persistentvolumeclaim/create [post]
func Create(c *gin.Context) {
	klog.Info("调用创建 PersistentVolumeClaim 对象的函数")

	var r CreatePersistentVolumeClaimRequest
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

	pvc := newPersistentVolumeClaim(r)
	result, err := clientset.CoreV1().PersistentVolumeClaims(r.Namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
	if err != nil {
		tool.SendResponse(c, errno.ErrCreatePersistentVolumeClaim, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

func newPersistentVolumeClaim(r CreatePersistentVolumeClaimRequest) *v1.PersistentVolumeClaim {
	pvc := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	spec := v1.PersistentVolumeClaimSpec{}

	if r.StorageClassName != nil {
		spec.StorageClassName = r.StorageClassName
	}

	if r.StorageCapacity > 0 {
		in := strconv.FormatFloat(r.StorageCapacity, 'f', 5, 32)
		capacity := in + viper.GetString("constants.storage_unit")
		klog.Infof("capacity is %s", capacity)
		request, _ := resource.ParseQuantity(capacity)
		spec.Resources = v1.ResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceStorage: request,
			},
		}
	}

	if len(r.AccessModes) > 0 {
		for _, accessMode := range r.AccessModes {
			var mode v1.PersistentVolumeAccessMode
			switch accessMode {
			case "ReadWriteOnce":
				mode = v1.ReadWriteOnce
			case "ReadOnlyMany":
				mode = v1.ReadOnlyMany
			case "ReadWriteMany":
				mode = v1.ReadWriteMany
			}
			spec.AccessModes = append(spec.AccessModes, mode)
		}
	}

	pvc.Spec = spec

	return &pvc
}
