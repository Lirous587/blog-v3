//go:build wireinject
// +build wireinject

package maxim

import (
	"blog/internal/domain/maxim/controller"
	"blog/internal/domain/maxim/infrastructure/cache"
	"blog/internal/domain/maxim/infrastructure/db"
	"blog/internal/domain/maxim/router"
	"blog/internal/domain/maxim/service"
	"blog/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) error {
	wire.Build(
		router.RegisterV1,
		controller.NewController,
		service.NewService,
		db.NewDB,
		cache.NewCache,
		repository.GormDB,
		repository.RedisClient,
	)
	return nil
}
