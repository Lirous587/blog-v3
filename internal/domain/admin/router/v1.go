package router

import (
	"blog/internal/domain/admin/controller"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, adminCtrl controller.Controller) error {
	g := r.Group("/v1/admin")
	{
		g.GET("/ifInit", adminCtrl.IfInit)
		g.POST("/init", adminCtrl.Init)
		g.POST("/login", adminCtrl.Login)
		g.POST("/refresh_token", adminCtrl.RefreshToken)
	}
	return nil
}
