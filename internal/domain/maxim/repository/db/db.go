package db

import (
	"blog/internal/domain/maxim/model"
	"blog/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DB interface {
	Create(admin *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	List(res *model.ListReq) (*model.ListRes, error)
	FindByID(id uint) (*model.Maxim, error)
	GetRandom20() ([]model.MaximDTO, error)
}

type db struct {
	orm *gorm.DB
}

func NewDB(orm *gorm.DB) DB {
	//var maxim model.Maxim
	//orm.AutoMigrate(&maxim)
	return &db{orm: orm}
}

func (db *db) Create(req *model.CreateReq) error {
	maxim := model.Maxim{
		Content: req.Content,
		Author:  req.Author,
		Color:   req.Color,
	}
	return db.orm.Create(&maxim).Error
}

func (db *db) Update(id uint, req *model.UpdateReq) error {
	maxim := model.Maxim{
		Content: req.Content,
		Author:  req.Author,
		Color:   req.Color,
	}
	return db.orm.Where("id = ?", id).Updates(&maxim).Error
}

func (db *db) Delete(id uint) error {
	return db.orm.Delete(&model.Maxim{}, id).Error
}

func (db *db) List(req *model.ListReq) (*model.ListRes, error) {
	keyword := utils.BuildLikeQuery(req.Keyword)

	var count int64
	if err := db.orm.Model(&model.Maxim{}).Where("content LIKE ? or author LIKE ?", keyword, keyword).Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	pages, err := utils.ComputePages(count, req.PageSize, req.Page)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	maxims := make([]model.Maxim, 0, req.PageSize)

	offset, err := utils.ComputeOffset(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.orm.Limit(req.PageSize).Offset(offset).Where("content LIKE ? or author LIKE ?", keyword, keyword).Find(&maxims).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	dtos := make([]model.MaximDTO, len(maxims))
	for i, maxim := range maxims {
		dtos[i] = *maxim.ConvertToDTO()
	}

	res := &model.ListRes{
		List:  dtos,
		Pages: pages,
	}

	return res, nil
}

func (db *db) FindByID(id uint) (*model.Maxim, error) {
	var maxim model.Maxim
	if err := db.orm.First(&maxim, id).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &maxim, nil
}

func (db *db) GetRandom20() ([]model.MaximDTO, error) {
	var maxims []model.Maxim
	randomFunc := utils.ResolveDBRandomFunc(db.orm)

	if err := db.orm.Order(randomFunc).Find(&maxims).Limit(20).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	dtos := make([]model.MaximDTO, len(maxims))

	for i := range maxims {
		dtos[i] = *maxims[i].ConvertToDTO()
	}

	return dtos, nil
}
