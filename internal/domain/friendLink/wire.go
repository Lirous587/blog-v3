//go:build wireinject
// +build wireinject

package friendLink

import (
	"blog/internal/domain/friendLink/controller"
	"blog/internal/domain/friendLink/infrastructure/cache"
	"blog/internal/domain/friendLink/infrastructure/db"
	"blog/internal/domain/friendLink/infrastructure/notifier"
	"blog/internal/domain/friendLink/router"
	"blog/internal/domain/friendLink/service"
	"blog/internal/domain/friendLink/worker"
	"blog/pkg/email"
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
		notifier.NewMailer,
		repository.GormDB,
		repository.RedisClient,
		email.GetMailer,
	)
	return nil
}

func InitWorker() worker.Worker {
	wire.Build(
		worker.NewWorker,
		db.NewDB,
		notifier.NewMailer,
		repository.GormDB,
		email.GetMailer,
	)
	return nil
}
