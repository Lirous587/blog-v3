package service

import (
	"blog/internal/domain/label/model"
	"blog/internal/domain/label/repository/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service interface {
	Create(req *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	List(req *model.ListReq) (res *model.ListRes, err error)
	FindByIDs(ids []uint) ([]model.Label, error)
}

type service struct {
	db db.DB
}

func NewService(db db.DB) Service {
	return &service{db: db}
}

func (s *service) Create(req *model.CreateReq) (err error) {
	// 先查是否有同名的记录
	label, err := s.db.FindByName(req.Name)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(err)
	}

	if label.Name == req.Name {
		return errors.New("label名重复")
	}

	if err = s.db.Create(req); err != nil {
		return errors.WithStack(err)
	}

	return
}

func (s *service) Update(id uint, req *model.UpdateReq) (err error) {
	// 除开当前id外 不得有同名的label
	_, err = s.db.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return errors.WithStack(err)
	}

	// 检查名称唯一性
	exists, err := s.db.IsNameTakenByOthers(req.Name, id)
	if err != nil {
		return errors.WithStack(err)
	}
	if exists {
		return errors.New("存在重复的label记录")
	}

	return errors.WithStack(s.db.Update(id, req))
}

func (s *service) Delete(id uint) (err error) {
	_, err = s.db.FindByID(id)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(s.db.Delete(id))
}

func (s *service) List(req *model.ListReq) (res *model.ListRes, err error) {
	res, err = s.db.List(req)
	return res, errors.WithStack(err)
}

func (s *service) FindByIDs(ids []uint) ([]model.Label, error) {
	return s.db.FindByIDs(ids)
}
