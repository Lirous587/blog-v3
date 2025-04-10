//go:build wireinject
// +build wireinject

package admin

import (
	"blog/internal/domain/admin/controller"
	"blog/internal/domain/admin/infrastructure/cache"
	"blog/internal/domain/admin/infrastructure/db"
	"blog/internal/domain/admin/router"
	"blog/internal/domain/admin/service"
	"blog/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) error {
	wire.Build(
		router.RegisterV1,
		controller.NewController,
		service.NewService,
		repository.GormDB,
		repository.RedisClient,
		db.NewDB,
		cache.NewCache,
	)
	return nil
}
