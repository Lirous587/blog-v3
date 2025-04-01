package controller

import (
	"blog/internal/domain/admin/model"
	"blog/internal/domain/admin/service"
	"blog/internal/response"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type controller struct {
	server service.Service
}

// NewAdminController 创建控制器实例
func NewController(service service.Service) Controller {
	return &controller{
		server: service,
	}
}

func (c *controller) Register(ctx *gin.Context) {
	req := new(model.RegisterReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ClientError(ctx, response.CodeParamInvalid, err)
		return
	}

	if err := c.server.Register(req); err != nil {
		response.ServerError(ctx, response.CodeServerError, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) Login(ctx *gin.Context) {
	req := new(model.LoginReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ClientError(ctx, response.CodeParamInvalid, err)
		return
	}

	res, err := c.server.Auth(req.Email, req.Password)
	if err != nil {
		response.ClientError(ctx, response.CodeAuthFailed, err)
		return
	}

	response.Success(ctx, res)
}
