package router

import (
	"blog/internal/domain/label/controller"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl controller.Controller) error {
	adminAuth, err := middleware.InitAdminAuth()
	if err != nil {
		panic(err)
	}

	authValidate := adminAuth.Validate()

	adminGroup := r.Group("/v1/label", authValidate)
	{
		adminGroup.POST("/", ctrl.Create)
		adminGroup.PUT("/:id", ctrl.Update)
		adminGroup.DELETE("/:id", ctrl.Delete)
		adminGroup.GET("/list", ctrl.List)
		adminGroup.GET("/all", ctrl.GetAll)
	}
	publicGroup := r.Group("/v1/label")
	publicGroup.GET("/all/with_essay_count", ctrl.GetAllWithEssayCount)

	return nil
}
