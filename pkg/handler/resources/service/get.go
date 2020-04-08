package service

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/service"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary  查询某一 Service 对象的详情
// @Description 查询某一 Service 对象的详情
// @Tags resource
// @Accept json
// @Produce json
// @param name path string true "Service 对象名称"
// @Param namespace path string true "用户的命名空间"
// @Success 200 {object} handler.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/service/detail/{name}/{namespace} [get]
func GetService(c *gin.Context) {
	log.Info("调用创建 Service 对象的函数")

	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		SendResponse(c, errno.ErrBadParam, nil)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	result, err := service.GetServiceDetail(clientset, namespace, name)
	if err != nil {
		SendResponse(c, errno.ErrGetService, err)
		return
	}

	SendResponse(c, errno.OK, result)
}
