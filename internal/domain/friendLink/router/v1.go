package router

import (
	"blog/internal/domain/friendLink/controller"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl controller.Controller) error {
	adminAuth, err := middleware.InitAdminAuth()
	if err != nil {
		panic(err)
	}

	authValidate := adminAuth.Validate()

	adminGroup := r.Group("/v1/friend_link", authValidate)
	{
		adminGroup.POST("/", ctrl.Create)
		adminGroup.PUT("/:id", ctrl.Update)
		adminGroup.PATCH("/:id/status", ctrl.UpdateStatus)
		adminGroup.DELETE("/:id", ctrl.Delete)
		adminGroup.GET("/list", ctrl.List)
	}

	publicGroup := r.Group("/v1/friend_link")
	{
		publicGroup.GET("/public/random20", ctrl.GetPublicRandom20)
		publicGroup.POST("/apply", ctrl.Apply)
	}
	return nil
}
