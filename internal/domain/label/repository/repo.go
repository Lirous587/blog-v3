package repository

import (
	"blog/internal/domain/label/model"
	"blog/utils"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	FindByName(name string) (*model.Label, error)
	FindByID(id uint) (*model.Label, error)
	IsNameTakenByOthers(name string, exceptID uint) (bool, error)
	Create(admin *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	List(res *model.ListReq) (*model.ListRes, error)
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
	if err := r.db.Where("name = ?", name).First(label).Error; err != nil {
		return nil, err
	}
	return label, nil
}

func (r *repository) FindByID(id uint) (*model.Label, error) {
	label := new(model.Label)
	if err := r.db.Where("id = ?", id).First(label).Error; err != nil {
		return nil, err
	}
	return label, nil
}

func (r *repository) IsNameTakenByOthers(name string, exceptID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Label{}).
		Where("name = ? AND id != ?", name, exceptID).
		Count(&count).Error

	return count > 0, err
}

func (r *repository) Create(req *model.CreateReq) error {
	label := model.Label{
		Name:         req.Name,
		Introduction: req.Introduction,
	}
	return r.db.Create(&label).Error
}

func (r *repository) Update(id uint, req *model.UpdateReq) error {
	label := model.Label{
		Name:         req.Name,
		Introduction: req.Introduction,
	}
	return r.db.Model(&model.Label{}).Where("id = ?", id).Updates(&label).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&model.Label{}, id).Error
}

func (r *repository) List(req *model.ListReq) (*model.ListRes, error) {
	keyword := utils.BuildLikeQuery(req.Keyword)

	var count int64
	if err := r.db.Model(&model.Label{}).Where("name LIKE ?", keyword).Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	pages := utils.ComputePages(count, req.PageSize)

	//  检查页码是否超出范围
	if req.Page > pages {
		return nil, errors.New("请求页码超出范围")
	}

	labels := make([]model.Label, 0, req.PageSize)
	offset := utils.ComputeOffset(req.Page, req.PageSize)

	if err := r.db.Limit(req.PageSize).Offset(offset).Where("name LIKE ?", keyword).Find(&labels).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	dtos := make([]model.LabelDTO, len(labels))
	for i, label := range labels {
		dtos[i] = model.LabelDTO{
			ID:           label.ID,
			Name:         label.Name,
			Introduction: label.Introduction,
			CreatedAt:    utils.FormatTime(label.CreatedAt),
		}
	}

	res := &model.ListRes{
		List:  dtos,
		Pages: pages,
	}

	return res, nil
}
