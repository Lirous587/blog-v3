//go:build wireinject
// +build wireinject

package moment

import (
	"blog/internal/domain/moment/controller"
	"blog/internal/domain/moment/infrastructure/cache"
	"blog/internal/domain/moment/infrastructure/db"
	"blog/internal/domain/moment/router"
	"blog/internal/domain/moment/service"
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
