//go:build wireinject
// +build wireinject

package label

import (
	"blog/internal/domain/label/controller"
	"blog/internal/domain/label/repository/cache"
	"blog/internal/domain/label/repository/db"
	"blog/internal/domain/label/router"
	"blog/internal/domain/label/service"
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
