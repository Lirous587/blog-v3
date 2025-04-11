package db

import (
	"blog/internal/domain/moment/model"
	"blog/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DB interface {
	Create(admin *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	List(res *model.ListReq) (*model.ListRes, error)
	FindByID(id uint) (*model.Moment, error)
	GetRandom20() ([]model.MomentDTO, error)
}

type db struct {
	orm *gorm.DB
}

func NewDB(orm *gorm.DB) DB {
	//var moment model.Moment
	//orm.AutoMigrate(&moment)
	return &db{orm: orm}
}

func (db *db) Create(req *model.CreateReq) error {
	moment := model.Moment{
		Title:    req.Title,
		Content:  req.Content,
		Location: req.Location,
	}
	return db.orm.Create(&moment).Error
}

func (db *db) Update(id uint, req *model.UpdateReq) error {
	moment := model.Moment{
		Title:    req.Title,
		Content:  req.Content,
		Location: req.Location,
	}
	return db.orm.Where("id = ?", id).Updates(&moment).Error
}

func (db *db) Delete(id uint) error {
	return db.orm.Delete(&model.Moment{}, id).Error
}

func (db *db) List(req *model.ListReq) (*model.ListRes, error) {
	keyword := utils.BuildLikeQuery(req.Keyword)

	var count int64
	if err := db.orm.Model(&model.Moment{}).Where("title LIKE ? or location LIKE ?", keyword, keyword).Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	pages, err := utils.ComputePages(count, req.PageSize, req.Page)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	moments := make([]model.Moment, 0, req.PageSize)

	offset, err := utils.ComputeOffset(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.orm.Limit(req.PageSize).Offset(offset).Where("title LIKE ? or location LIKE ?", keyword, keyword).Find(&moments).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	dtos := make([]model.MomentDTO, len(moments))
	for i, moment := range moments {
		dtos[i] = *moment.ConvertToDTO()
	}

	res := &model.ListRes{
		List:  dtos,
		Pages: pages,
	}

	return res, nil
}

func (db *db) FindByID(id uint) (*model.Moment, error) {
	var moment model.Moment
	if err := db.orm.First(&moment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return &moment, nil
}

func (db *db) GetRandom20() ([]model.MomentDTO, error) {
	var moments []model.Moment
	randomFunc := utils.ResolveDBRandomFunc(db.orm)

	if err := db.orm.Order(randomFunc).Find(&moments).Limit(20).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	dtos := make([]model.MomentDTO, len(moments))

	for i := range moments {
		dtos[i] = *moments[i].ConvertToDTO()
	}

	return dtos, nil
}
