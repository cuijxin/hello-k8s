package user

import (
	"hello-k8s/pkg/errno"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func Create(c *gin.Context) {
	log.Debug("调用创建用户的接口！")
	var r struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var err error
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, err)
		return
	}

    log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
    if r.Username == "" {
        err = errno.New(errno.ErrUserNotFound, fmt.)
    }
}
