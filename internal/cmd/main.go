package cmd

import (
	"blog/internal/domain/admin"
	"blog/internal/domain/label"
	"blog/pkg/httpserver"
)

func Main() {
	// 创建服务器
	s := httpserver.New(8080)
	r := s.Router

	// 创建 /api 分组
	api := r.Group("/api")

	var err error
	// 创建admin路由
	if err = admin.InitV1(api); err != nil {
		panic(err)
	}

	if err := label.InitV1(api); err != nil {
		panic(err)
	}

	//api.GET("/auth", adminAuth.Validate(), func(ctx *gin.Context) {
	//	ctx.JSON(200, gin.H{
	//		"msg": "认证成功",
	//	})
	//})
	s.Run()
}
