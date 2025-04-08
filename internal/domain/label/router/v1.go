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

	g := r.Group("/v1/label", authValidate)
	{
		g.POST("/", ctrl.Create)
		g.PUT("/:id", ctrl.Update)
		g.DELETE("/:id", ctrl.Delete)
		g.GET("/list", ctrl.List)
	}
	g2 := r.Group("/v1/label")
	g2.GET("/all", ctrl.GetAllWithEssayCount)

	return nil
}
