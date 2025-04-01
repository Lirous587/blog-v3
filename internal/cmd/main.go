package cmd

import (
	"blog/internal/domain/admin"
	"blog/internal/middleware"
	"blog/pkg/httpserver"
)

func Main() {
	// 创建服务器
	s := httpserver.New(8080)
	r := s.Router

	// 创建 /api 分组
	api := r.Group("/api")
	api.Use(middleware.ErrorHandler())

	// 创建admin路由
	if err := admin.InitV1(api); err != nil {
		panic(err)
	}

	s.Run()
}
