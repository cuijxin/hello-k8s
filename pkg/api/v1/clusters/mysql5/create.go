package mysql5

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
)

// @Summary 创建MySQL V5集群
// @Description 创建MySQL V5集群
// @Tags add-on
// @Accept json
// @Produce json
// @param data body mysql5.ClusterOptions true "创建 MySQL V5集群所需参数."
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/addon/mysql5/cluster/create [post]
func CreateCluster(c *gin.Context) {
	klog.Info("调用创建 MySQL V5 集群函数")

	var r ClusterOptions
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
	}
	klog.Infof("request body is [%v]", r)

	tool.SendResponse(c, errno.OK, nil)

	return
}
