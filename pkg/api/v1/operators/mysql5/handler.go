package mysql5

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/component/addons"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
)

// @Summary 在Kubernetes集群中安装 MySQL V5 Operator.
// @Description 在Kubernetes集群中安装 MySQL V5 Operator.
// @Tags add-on
// @Accept json
// @Produce json
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/addon/mysql5/operator [post]
func InstallOperator(c *gin.Context) {
	klog.Info("调用安装 MySQL V5 版本的 Operator 的函数.")

	mysql5Operator := New()

	if err := mysql5Operator.Deploy(client.MyClient, addons.AddOnOptions{}); err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
	return
}

// @Summary 从Kubernetes集群中删除MySQL V5 Operator.
// @Description 从Kubernetes集群中删除MySQL V5 Operator.
// @Tags add-on
// @Accept json
// @Produce json
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/addon/mysql5/operator [delete]
func UnInstallOperator(c *gin.Context) {
	klog.Info("调用卸载 MySQL V5 版本的 Operator 的函数.")

	mysql5Operator := New()

	if err := mysql5Operator.Delete(client.MyClient); err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
	return
}
