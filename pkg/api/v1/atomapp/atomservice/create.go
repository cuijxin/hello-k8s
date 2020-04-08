package atomservice

import (
	"hello-k8s/pkg/api/v1/tool"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"

	"hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary 创建Atom自定义服务.
// @Description 创建Atom自定义服务.
// @Tags atomapp
// @Accept json
// @Produce json
// @param data body deployment.AppDeploymentSpec true "创建Atom自定义服务时所需参数."
// @Success 200 {object} atomservice.CreateAtomServiceResponse "{"code":0,"message":"OK","data":{""}}"
// @Router /atomapp/atomservice/create [post]
func Create(c *gin.Context) {
	log.Info("调用创建Atom自定义服务的函数.")

	var r *deployment.AppDeploymentSpec
	if err := c.BindJSON(r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	CreateNamespace(r.Namespace, clientset)

	if err := deployment.DeployApp(r, clientset); err != nil {
		tool.SendResponse(c, errno.ErrDeployAtomService, err)
		return
	}

	rsp := CreateAtomServiceResponse{
		Namespace:  r.Namespace,
		Deployment: r.Name,
		Service:    r.Name,
	}

	tool.SendResponse(c, errno.OK, rsp)
}
