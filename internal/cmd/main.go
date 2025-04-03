package cmd

import (
	"blog/internal/domain/admin"
	"blog/internal/middleware"
	"blog/pkg/httpserver"

	"github.com/gin-gonic/gin"
)

func Main() {
	// 创建服务器
	s := httpserver.New(8080)
	r := s.Router

	// 创建 /api 分组
	api := r.Group("/api")
	api.Use(middleware.ErrorHandler())

	adminAuth, err := middleware.InitAdminAuth()
	if err != nil {
		panic(err)
	}

	// 创建admin路由
	if err = admin.InitV1(api); err != nil {
		panic(err)
	}

	api.GET("/auth", adminAuth.Validate(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"msg": "认证成功",
		})
	})

	s.Run()
}
