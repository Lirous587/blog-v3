//go:build wireinject
// +build wireinject

package admin

import (
	"blog/internal/domain/admin/controller"
	"blog/internal/domain/admin/repository"
	"blog/internal/domain/admin/router"
	"blog/internal/domain/admin/service"
	"blog/pkg/repository/db"
	"blog/pkg/repository/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// InitializeAdminAPI 初始化Admin模块的API
func InitV1(r *gin.RouterGroup) error {
	wire.Build(
		router.RegisterV1,
		controller.NewController,
		service.NewService,
		repository.NewRepository,
		db.DB,
		redis.Client,
	)
	return nil
}
