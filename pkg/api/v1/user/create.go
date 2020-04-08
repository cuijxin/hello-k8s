package user

import (
	model "hello-k8s/pkg/model/user"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// @Summary 创建 User 对象
// @Description 创建 User 对象
// @Tags user
// @Accept json
// @Produce json
// @param data body CreateRequest true "创建 User 对象时所需参数"
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/user [post]
func Create(c *gin.Context) {
	log.Debug("调用创建用户的接口！", lager.Data{"X-Request-Id": tool.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, nil)
		return
	}

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// Validate the data.
	if err := u.Validate(); err != nil {
		tool.SendResponse(c, errno.ErrValidation, nil)
		return
	}
	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		tool.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}
	// Insert the user to the database.
	if err := u.Create(); err != nil {
		tool.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	tool.SendResponse(c, nil, rsp)
}
