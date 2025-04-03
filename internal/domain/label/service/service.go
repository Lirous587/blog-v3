package service

import (
	"blog/internal/domain/label/model"
	"blog/internal/domain/label/repository"
	"github.com/pkg/errors"
)

type Service interface {
	Create(req *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint, req *model.DeleteReq) error
	List(req *model.ListReq) (res model.ListRes, err error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s service) Create(req *model.CreateReq) (err error) {
	label, err := s.repo.FindByName(req.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	if label != nil {
		return errors.New("数据已存在")
	}

	if err = s.repo.Create(req); err != nil {
		return
	}

	return
}

func (s service) Update(id uint, req *model.UpdateReq) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Delete(id uint, req *model.DeleteReq) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s service) List(req *model.ListReq) (res model.ListRes, err error) {
	//TODO implement me
	panic("implement me")
}
