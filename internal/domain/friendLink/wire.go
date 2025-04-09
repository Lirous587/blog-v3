//go:build wireinject
// +build wireinject

package friendLink

import (
	"blog/internal/domain/friendLink/controller"
	"blog/internal/domain/friendLink/repository/cache"
	"blog/internal/domain/friendLink/repository/db"
	"blog/internal/domain/friendLink/router"
	"blog/internal/domain/friendLink/service"
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
