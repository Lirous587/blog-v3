package controller

import (
	"blog/internal/domain/essay/model"
	"blog/internal/domain/essay/service"
	"blog/pkg/response"
	"blog/utils"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	List(ctx *gin.Context)
	GetTimelines(ctx *gin.Context)
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
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	if err := c.server.Create(req); err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) Update(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	req := new(model.UpdateReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	if err := c.server.Update(id, req); err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}

	response.Success(ctx)
}

func (c *controller) Delete(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	if err = c.server.Delete(id); err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) Get(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}
	res, err := c.server.Get(id)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}

	response.Success(ctx, res)
}

func (c *controller) List(ctx *gin.Context) {
	req := new(model.ListReq)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	res, err := c.server.List(req)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}

	response.Success(ctx, res)
}

func (c *controller) GetTimelines(ctx *gin.Context) {
	req := new(model.TimelineReq)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.Error(ctx, response.CodeParamInvalid, err)
		return
	}

	res, err := c.server.GetTimelines(req)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err)
		return
	}

	response.Success(ctx, res)
}
