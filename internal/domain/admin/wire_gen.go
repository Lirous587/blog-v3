// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package admin

import (
	"blog/internal/domain/admin/controller"
	"blog/internal/domain/admin/repository/cache"
	"blog/internal/domain/admin/repository/db"
	"blog/internal/domain/admin/router"
	"blog/internal/domain/admin/service"
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
