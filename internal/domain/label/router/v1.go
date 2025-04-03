package router

import (
	"blog/internal/domain/label/controller"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, LabelCtrl controller.Controller) error {
	adminAuth, err := middleware.InitAdminAuth()
	if err != nil {
		panic(err)
	}

	authValidate := adminAuth.Validate()

	g := r.Group("/v1/label")
	{
		g.POST("/", authValidate, LabelCtrl.Create)
		g.PUT("/:id", authValidate, LabelCtrl.Update)
		g.DELETE("/:id", authValidate, LabelCtrl.Delete)
		g.GET("/list", LabelCtrl.List)
	}
	return nil
}
