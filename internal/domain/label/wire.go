//go:build wireinject
// +build wireinject

package label

import (
	"blog/internal/domain/label/controller"
	"blog/internal/domain/label/repository"
	"blog/internal/domain/label/router"
	"blog/internal/domain/label/service"
	"blog/pkg/repository/db"
	"blog/pkg/repository/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

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
