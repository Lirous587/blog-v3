package cmd

import (
	"blog/internal/domain/admin"
	"blog/internal/domain/essay"
	"blog/internal/domain/friendLink"
	"blog/internal/domain/label"
	"blog/internal/domain/maxim"
	"blog/pkg/httpserver"
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
	{
		essayWorker := essay.InitWorker()
		essayWorker.Start()
		s.RegisterStopHandler(func() {
			essayWorker.Stop()
		})
	}

	if err := maxim.InitV1(api); err != nil {
		panic(err)
	}

	if err := friendLink.InitV1(api); err != nil {
		panic(err)
	}

	s.Run()
}
