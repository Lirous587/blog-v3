package controller

import (
	"blog/internal/domain/label/model"
	"blog/internal/domain/label/service"
	"blog/internal/response"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	List(ctx *gin.Context)
}

func NewController(service service.Service) Controller {
	return &controller{
		server: service,
	}
}

func (c *controller) Create(ctx *gin.Context) {
	req := new(model.CreateReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	if err := c.server.Create(req); err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) Delete(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c *controller) Update(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c *controller) List(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

type controller struct {
	server service.Service
}
