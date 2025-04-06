package router

import (
	"blog/internal/domain/essay/controller"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl controller.Controller) error {
	adminAuth, err := middleware.InitAdminAuth()
	if err != nil {
		panic(err)
	}

	authValidate := adminAuth.Validate()

	g := r.Group("/v1/essay")
	{
		g.POST("/", authValidate, ctrl.Create)
		g.PUT("/:id", authValidate, ctrl.Update)
		g.DELETE("/:id", authValidate, ctrl.Delete)
		g.GET("/:id", ctrl.Get)
		g.GET("/list", ctrl.List)
		g.GET("/timeline", ctrl.GetTimelines)
	}

	return nil
}
