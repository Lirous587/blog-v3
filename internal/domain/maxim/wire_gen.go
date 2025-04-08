// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package maxim

import (
	"blog/internal/domain/maxim/controller"
	"blog/internal/domain/maxim/repository/cache"
	"blog/internal/domain/maxim/repository/db"
	"blog/internal/domain/maxim/router"
	"blog/internal/domain/maxim/service"
	"blog/pkg/repository"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitV1(r *gin.RouterGroup) error {
	gormDB := repository.GormDB()
	dbDB := db.NewDB(gormDB)
	client := repository.RedisClient()
	cacheCache := cache.NewCache(client)
	serviceService := service.NewService(dbDB, cacheCache)
	controllerController := controller.NewController(serviceService)
	error2 := router.RegisterV1(r, controllerController)
	return error2
}
