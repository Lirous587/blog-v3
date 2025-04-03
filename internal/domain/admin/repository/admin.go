package repository

import (
	"blog/internal/domain/admin/model"
	"blog/pkg/config"
	"blog/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	HaveOne() (bool, error)
	FindByEmail(email string) (*model.Admin, error)
	Create(admin *model.Admin) error
	GenRefreshToken(payload *model.JwtPayload) (string, error)
	ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error
}

type repository struct {
	db     *gorm.DB
	client *redis.Client
}

func NewRepository(db *gorm.DB, client *redis.Client) Repository {
	// var admin model.Admin
	// db.AutoMigrate(&admin)
	return &repository{db: db, client: client}
}

func (r *repository) HaveOne() (bool, error) {
	var admin model.Admin
	err := r.db.First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.WithStack(err)
	}
	if &admin != nil {
		return true, nil
	}
	return false, nil
}

func (r *repository) FindByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Model(admin).Where("email = ?", email).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return &admin, nil
}

func (r *repository) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

const (
	keyRefreshToken = "refreshToken"
)

func (r *repository) GenRefreshToken(payload *model.JwtPayload) (string, error) {
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

	pipe := r.client.Pipeline()

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

func (r *repository) ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error {
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	key := utils.GetRedisKey(keyRefreshToken)

	result, err := r.client.HGet(context.Background(), key, string(payloadByte)).Result()
	if err != nil {
		return errors.WithStack(err)
	}

	if refreshToken == result {
		return nil
	}
	return errors.New("refreshToken 无效")
}
