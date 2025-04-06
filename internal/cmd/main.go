package cmd

import (
	"blog/internal/domain/admin"
	"blog/internal/domain/essay"
	"blog/internal/domain/label"
	"blog/pkg/httpserver"
	"go.uber.org/zap"
)

func Main() {
	// 创建服务器
	s := httpserver.New(8080)
	r := s.Router

	// 创建 /api 分组
	api := r.Group("/api")

	var err error

	if err = admin.InitV1(api); err != nil {
		panic(err)
	}

	if err := label.InitV1(api); err != nil {
		panic(err)
	}

	if err := essay.InitV1(api); err != nil {
		panic(err)
	}
	essayWorker := essay.InitWorker()
	essayWorker.Start()
	// 注册worker关闭函数
	s.RegisterStopHandler(func() {
		zap.L().Info("正在关闭Worker...")
		essayWorker.Stop()
	})

	s.Run()
}
