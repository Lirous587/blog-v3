package db

import (
	"blog/internal/domain/label/model"
	"blog/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DB interface {
	FindByName(name string) (*model.Label, error)
	FindByID(id uint) (*model.Label, error)
	FindByIDs(ids []uint) ([]model.Label, error)
	IsNameTakenByOthers(name string, exceptID uint) (bool, error)
	Create(admin *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	List(res *model.ListReq) (*model.ListRes, error)
}

type db struct {
	orm *gorm.DB
}

func NewDB(orm *gorm.DB) DB {
	//var label model.Label
	//orm.AutoMigrate(&label)
	return &db{orm: orm}
}

func (db *db) FindByName(name string) (*model.Label, error) {
	label := new(model.Label)
	if err := db.orm.Where("name = ?", name).First(label).Error; err != nil {
		return label, err
	}
	return label, nil
}

func (db *db) FindByID(id uint) (*model.Label, error) {
	label := new(model.Label)
	if err := db.orm.Where("id = ?", id).First(label).Error; err != nil {
		return nil, err
	}
	return label, nil
}

func (db *db) FindByIDs(ids []uint) ([]model.Label, error) {
	var labels []model.Label
	if len(ids) == 0 {
		return labels, errors.New("无效的ids输入")
	}
	if err := db.orm.Where("id IN ?", ids).Find(&labels).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return labels, nil
}

func (db *db) IsNameTakenByOthers(name string, exceptID uint) (bool, error) {
	var count int64
	err := db.orm.Model(&model.Label{}).
		Where("name = ? AND id != ?", name, exceptID).
		Count(&count).Error

	return count > 0, err
}

func (db *db) Create(req *model.CreateReq) error {
	label := model.Label{
		Name:         req.Name,
		Introduction: req.Introduction,
	}
	return db.orm.Create(&label).Error
}

func (db *db) Update(id uint, req *model.UpdateReq) error {
	label := model.Label{
		Name:         req.Name,
		Introduction: req.Introduction,
	}
	return db.orm.Model(&model.Label{}).Where("id = ?", id).Updates(&label).Error
}

func (db *db) Delete(id uint) error {
	return db.orm.Delete(&model.Label{}, id).Error
}

func (db *db) List(req *model.ListReq) (*model.ListRes, error) {
	keyword := utils.BuildLikeQuery(req.Keyword)

	var count int64
	if err := db.orm.Model(&model.Label{}).Where("name LIKE ?", keyword).Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	pages, err := utils.ComputePages(count, req.PageSize, req.Page)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	labels := make([]model.Label, 0, req.PageSize)

	offset, err := utils.ComputeOffset(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.orm.Limit(req.PageSize).Offset(offset).Where("name LIKE ?", keyword).Find(&labels).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	dtos := make([]model.LabelDTO, len(labels))
	for i, label := range labels {
		dtos[i] = *label.ConvertToDTO()
	}

	res := &model.ListRes{
		List:  dtos,
		Pages: pages,
	}

	return res, nil
}
