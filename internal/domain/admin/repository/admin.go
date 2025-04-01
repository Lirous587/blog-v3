package repository

import (
	"blog/internal/domain/admin/model"
	"blog/utils"
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*model.Admin, error)
	Create(admin *model.Admin) error
	GenRefreshToken(email string) (string, error)
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

func (r *repository) FindByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
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

func (r *repository) GenRefreshToken(email string) (string, error) {
	refreshToken, err := utils.GenRandomHexToken()
	if err != nil {
		return "", errors.WithStack(err)
	}
	key := utils.GetRedisKey(keyRefreshToken)
	if err := r.client.HSet(context.Background(), key, email, refreshToken).Err(); err != nil {
		return "", errors.WithStack(err)
	}

	return refreshToken, nil
}
