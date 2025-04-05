package cache

import (
	"blog/internal/domain/admin/model"
	"blog/pkg/config"
	"blog/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GenRefreshToken(payload *model.JwtPayload) (string, error)
	ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	// var admin model.Admin
	// db.AutoMigrate(&admin)
	return &cache{client: client}
}

const (
	keyRefreshToken = "refreshToken"
)

func (ch *cache) GenRefreshToken(payload *model.JwtPayload) (string, error) {
	refreshToken, err := utils.GenRandomHexToken()
	if err != nil {
		return "", errors.WithStack(err)
	}
	key := utils.GetRedisKey(keyRefreshToken)
	payloadByte, err := json.Marshal(payload)
	payloadStr := string(payloadByte)
	if err != nil {
		return "", errors.WithStack(err)
	}

	pipe := ch.client.Pipeline()

	if err := pipe.HSet(context.Background(), key, payloadStr, refreshToken).Err(); err != nil {
		return "", errors.WithStack(err)
	}
	refreshExpireDuration := time.Duration(config.Cfg.Auth.RefreshToken.ExpireDay) * 24 * time.Hour

	// 使用计算出的过期时间
	pipe.HExpire(context.Background(), key, refreshExpireDuration, payloadStr)

	// 执行Pipeline命令
	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", errors.WithStack(err)
	}

	return refreshToken, nil
}

func (ch *cache) ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error {
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	key := utils.GetRedisKey(keyRefreshToken)

	result, err := ch.client.HGet(context.Background(), key, string(payloadByte)).Result()
	if err != nil {
		return errors.WithStack(err)
	}

	if refreshToken == result {
		return nil
	}
	return errors.New("refreshToken 无效")
}
