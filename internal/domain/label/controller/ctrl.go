package controller

import (
	"blog/internal/domain/label/model"
	"blog/internal/domain/label/service"
	"blog/pkg/response"
	"blog/utils"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	List(ctx *gin.Context)
	GetAllWithEssayCount(ctx *gin.Context)
}

type controller struct {
	server service.Service
}

func NewController(service service.Service) Controller {
	return &controller{
		server: service,
	}
}

func (c *controller) Create(ctx *gin.Context) {
	req := new(model.CreateReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	if err := c.server.Create(req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) Update(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	req := new(model.UpdateReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	if err := c.server.Update(id, req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

func (c *controller) Delete(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	if err := c.server.Delete(id); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) List(ctx *gin.Context) {
	req := new(model.ListReq)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, err := c.server.List(req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, res)
}

func (c *controller) GetAllWithEssayCount(ctx *gin.Context) {
	res, err := c.server.GetAllWithEssayCount()
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, res)
}
