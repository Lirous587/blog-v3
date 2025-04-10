package db

import (
	"blog/internal/domain/friendLink/model"
	"blog/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DB interface {
	Create(admin *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	UpdateStatus(id uint, req *model.UpdateStatusReq) error
	Delete(id uint) error
	List(res *model.ListReq) (*model.ListRes, error)
	FindByID(id uint) (*model.FriendLink, error)
	FindByURL(url string) (*model.FriendLink, error)
	GetPublishedRandom20() ([]model.FriendLinkDTO, error)
	Apply(req *model.ApplyReq) error
	GetPendingLinks() ([]model.FriendLink, error)
}

type db struct {
	orm *gorm.DB
}

func NewDB(orm *gorm.DB) DB {
	//var friendLink model.FriendLink
	//orm.AutoMigrate(&friendLink)
	return &db{orm: orm}
}

func (db *db) Create(req *model.CreateReq) error {
	friendLink := model.FriendLink{
		Introduction: req.Introduction,
		SiteName:     req.SiteName,
		Url:          req.Url,
		Logo:         req.Logo,
		Email:        req.Email,
	}
	return db.orm.Create(&friendLink).Error
}

func (db *db) Update(id uint, req *model.UpdateReq) error {
	friendLink := model.FriendLink{
		Introduction: req.Introduction,
		SiteName:     req.SiteName,
		Url:          req.Url,
		Logo:         req.Logo,
		Status:       model.StatusPending,
		Email:        req.Email,
	}
	return db.orm.Where("id = ?", id).Updates(&friendLink).Error
}

func (db *db) UpdateStatus(id uint, req *model.UpdateStatusReq) error {
	friendLink := model.FriendLink{
		Status: req.Status,
	}
	return db.orm.Where("id = ?", id).Updates(&friendLink).Error
}

func (db *db) Delete(id uint) error {
	return db.orm.Unscoped().Delete(&model.FriendLink{}, id).Error
}

func (db *db) List(req *model.ListReq) (*model.ListRes, error) {
	keyword := utils.BuildLikeQuery(req.Keyword)

	var count int64
	if err := db.orm.Model(&model.FriendLink{}).Where("status = ? AND site_name LIKE ?", model.StatusPublished, keyword).Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	pages, err := utils.ComputePages(count, req.PageSize, req.Page)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	friendLinks := make([]model.FriendLink, 0, req.PageSize)

	offset, err := utils.ComputeOffset(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.orm.Limit(req.PageSize).Offset(offset).Where("status = ? AND site_name LIKE ?", model.StatusPublished, keyword).Find(&friendLinks).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	dtos := make([]model.FriendLinkDTO, len(friendLinks))
	for i, friendLink := range friendLinks {
		dtos[i] = *friendLink.ConvertToDTO()
	}

	res := &model.ListRes{
		List:  dtos,
		Pages: pages,
	}

	return res, nil
}

func (db *db) FindByID(id uint) (*model.FriendLink, error) {
	var friendLink model.FriendLink
	if err := db.orm.First(&friendLink, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return &friendLink, nil
}

func (db *db) FindByURL(url string) (*model.FriendLink, error) {
	var friendLink model.FriendLink
	if err := db.orm.Where("url = ?", url).First(&friendLink).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return &friendLink, nil
}

func (db *db) GetPublishedRandom20() ([]model.FriendLinkDTO, error) {
	var maxims []model.FriendLink
	randomFunc := utils.ResolveDBRandomFunc(db.orm)

	if err := db.orm.Where("status = ?", model.StatusPublished).Order(randomFunc).Find(&maxims).Limit(20).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	dtos := make([]model.FriendLinkDTO, len(maxims))

	for i := range maxims {
		dtos[i] = *maxims[i].ConvertToDTO()
	}

	return dtos, nil
}

func (db *db) Apply(req *model.ApplyReq) error {
	friendLink := model.FriendLink{
		Introduction: req.Introduction,
		SiteName:     req.SiteName,
		Url:          req.Url,
		Logo:         req.Logo,
		Email:        req.Email,
	}
	return db.orm.Create(&friendLink).Error
}

func (db *db) GetPendingLinks() ([]model.FriendLink, error) {
	var list []model.FriendLink

	if err := db.orm.Where("status = ?", model.StatusPending).Find(&list).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}
