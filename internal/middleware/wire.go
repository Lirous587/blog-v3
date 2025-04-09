//go:build wireinject
// +build wireinject

package middleware

import (
	"blog/pkg/config"
	"github.com/google/wire"
)

// InitializeAdminAPI 初始化Admin模块的API
func InitAdminAuth() (Auth, error) {
	wire.Build(
		provideJWTSecret,
		NewAdminAuth,
	)
	return nil, nil
}

// provideJWTSecret 提供JWT Secret
func provideJWTSecret() []byte {
	return []byte(config.Cfg.JWT.Secret)
}
