// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package label

import (
	"blog/internal/domain/label/controller"
	"blog/internal/domain/label/repository"
	"blog/internal/domain/label/router"
	"blog/internal/domain/label/service"
	"blog/pkg/repository/db"
	"blog/pkg/repository/redis"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitV1(r *gin.RouterGroup) error {
	gormDB := db.DB()
	client := redis.Client()
	repositoryRepository := repository.NewRepository(gormDB, client)
	serviceService := service.NewService(repositoryRepository)
	controllerController := controller.NewController(serviceService)
	error2 := router.RegisterV1(r, controllerController)
	return error2
}
