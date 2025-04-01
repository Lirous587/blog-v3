package router

import (
	"blog/internal/domain/admin/controller"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, adminCtrl controller.Controller) error {
	g := r.Group("/v1/admin")
	{
		g.POST("/register", adminCtrl.Register)
		g.POST("/login", adminCtrl.Login)
	}
	return nil
}
