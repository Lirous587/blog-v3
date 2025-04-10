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
		adminGroup.DELETE("/:id", ctrl.Delete)
		adminGroup.GET("/list", ctrl.List)

		adminGroup.PATCH("/:id/approve", ctrl.Approve)
		adminGroup.DELETE("/:id/reject", ctrl.Reject)
	}

	publicGroup := r.Group("/v1/friend_link")
	{
		publicGroup.GET("/published/random20", ctrl.GetPublishedRandom20)
		publicGroup.POST("/apply", ctrl.Apply)
	}
	return nil
}
