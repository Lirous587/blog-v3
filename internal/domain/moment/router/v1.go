package router

import (
	"blog/internal/domain/moment/controller"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl controller.Controller) error {
	adminAuth, err := middleware.InitAdminAuth()
	if err != nil {
		panic(err)
	}

	authValidate := adminAuth.Validate()

	g := r.Group("/v1/moment", authValidate)
	{
		g.POST("/", ctrl.Create)
		g.PUT("/:id", ctrl.Update)
		g.DELETE("/:id", ctrl.Delete)
		g.GET("/list", ctrl.List)
		g.GET("/random20", ctrl.GetRandom20)
	}
	return nil
}
