package controller

import (
	"blog/internal/domain/admin/model"
	"blog/internal/domain/admin/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	IfInit(ctx *gin.Context)
	Init(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type controller struct {
	server service.Service
}

func NewController(service service.Service) Controller {
	return &controller{
		server: service,
	}
}

func (c *controller) IfInit(ctx *gin.Context) {
	have, err := c.server.IfInit()
	if err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}
	if have {
		response.Error(ctx, response.CodeAdminExist, nil)
		return
	}

	response.Success(ctx)
}

func (c *controller) Init(ctx *gin.Context) {
	req := new(model.InitReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}
	appErr := c.server.Init(req)
	if appErr != nil {
		response.ErrorStrict(ctx, appErr)
		return
	}
	response.Success(ctx)
}

func (c *controller) Login(ctx *gin.Context) {
	req := new(model.LoginReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	res, appErr := c.server.Auth(req.Email, req.Password)
	if appErr != nil {
		response.ErrorStrict(ctx, appErr)
		return
	}

	response.Success(ctx, res)
}

func (c *controller) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("refresh-token")
	if refreshToken == "" {
		response.Error(ctx, response.CodeAuthFailed, nil)
		return
	}

	req := new(model.RefreshTokenReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	res, appErr := c.server.RefreshToken(req, refreshToken)
	if appErr != nil {
		response.ErrorStrict(ctx, appErr)
		return
	}

	response.Success(ctx, res)
}
