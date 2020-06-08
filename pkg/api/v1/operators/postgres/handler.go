package postgres

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/component/addons"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
)

// @Summary 在Kubernetes集群中安装 Postgres Operator.
// @Description 在Kubernetes集群中安装 Postgres Operator.
// @Tags add-on
// @Accept json
// @Produce json
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/addon/postgres/operator [post]
func InstallOperator(c *gin.Context) {
	klog.Info("调用安装 Postgres Operator 的函数.")

	postgresOperator := New()

	if err := postgresOperator.Deploy(client.MyClient, addons.AddOnOptions{}); err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
	return
}

// @Summary 从Kubernetes集群中删除Postgres Operator.
// @Description 从Kubernetes集群中删除Postgres Operator.
// @Tags add-on
// @Accept json
// @Produce json
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/addon/postgres/operator [delete]
func UnInstallOperator(c *gin.Context) {
	klog.Info("调用卸载 Postgres Operator 的函数.")

	postgresOperator := New()

	if err := postgresOperator.Delete(); err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
	return
}
