package repository

import (
	"blog/internal/domain/label/model"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	FindByName(name string) (*model.Label, error)
	Create(admin *model.CreateReq) error
}

type repository struct {
	db     *gorm.DB
	client *redis.Client
}

func NewRepository(db *gorm.DB, client *redis.Client) Repository {
	//var Label model.Label
	//db.AutoMigrate(&Label)
	return &repository{db: db, client: client}
}

func (r *repository) FindByName(name string) (*model.Label, error) {
	label := new(model.Label)
	if err := r.db.Model(&model.Label{}).Where("name = ?", name).Find(label).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return nil, nil
}

func (r *repository) Create(req *model.CreateReq) error {
	label := model.Label{
		Name:         req.Name,
		Introduction: req.Introduction,
	}
	return r.db.Create(&label).Error
}

const (
	keyRefreshToken = "refreshToken"
)
