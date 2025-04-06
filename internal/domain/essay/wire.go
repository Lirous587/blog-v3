//go:build wireinject
// +build wireinject

package essay

import (
	"blog/internal/domain/essay/controller"
	essayCache "blog/internal/domain/essay/repository/cache"
	labelDB "blog/internal/domain/essay/repository/db"
	"blog/internal/domain/essay/router"
	"blog/internal/domain/essay/service"
	"blog/internal/domain/essay/worker"
	essayDB "blog/internal/domain/label/repository/db"
	labelServer "blog/internal/domain/label/service"
	"blog/pkg/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// 数据层依赖
var dataSet = wire.NewSet(
	repository.GormDB,
	repository.RedisClient,
)

// 标签领域依赖
var labelSet = wire.NewSet(
	labelDB.NewDB,
	labelServer.NewService,
)

// 文章领域依赖
var essaySet = wire.NewSet(
	essayDB.NewDB,
	essayCache.NewCache,
	service.NewService,
)

var routerSet = wire.NewSet(
	controller.NewController,
	router.RegisterV1,
)

func InitV1(r *gin.RouterGroup) error {
	wire.Build(
		routerSet,
		dataSet,
		labelSet,
		essaySet,
	)
	return nil
}

func InitWorker() worker.Worker {
	wire.Build(
		worker.NewWorker,
		dataSet,
		labelSet,
		essaySet,
	)
	return nil
}
