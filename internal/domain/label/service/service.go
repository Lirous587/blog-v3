package service

import (
	"blog/internal/domain/label/model"
	"blog/internal/domain/label/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service interface {
	Create(req *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	List(req *model.ListReq) (res *model.ListRes, err error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req *model.CreateReq) (err error) {
	// 先查是否有同名的记录
	_, err = s.repo.FindByName(req.Name)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(err)
	}

	if err = s.repo.Create(req); err != nil {
		return errors.WithStack(err)
	}

	return
}

func (s *service) Update(id uint, req *model.UpdateReq) (err error) {
	// 除开当前id外 不得有同名的label
	_, err = s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return errors.WithStack(err)
	}

	// 检查名称唯一性
	exists, err := s.repo.IsNameTakenByOthers(req.Name, id)
	if err != nil {
		return errors.WithStack(err)
	}
	if exists {
		return errors.New("存在重复的label记录")
	}

	return errors.WithStack(s.repo.Update(id, req))
}

func (s *service) Delete(id uint) (err error) {
	_, err = s.repo.FindByID(id)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(s.repo.Delete(id))
}

func (s *service) List(req *model.ListReq) (res *model.ListRes, err error) {
	res, err = s.repo.List(req)
	return res, errors.WithStack(err)
}
