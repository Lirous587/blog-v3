package service

import (
	"blog/internal/domain/maxim/model"
	"blog/internal/domain/maxim/repository/cache"
	"blog/internal/domain/maxim/repository/db"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Create(req *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	List(req *model.ListReq) (res *model.ListRes, err error)
	GetRandom20() ([]model.MaximDTO, error)
}

type service struct {
	db    db.DB
	cache cache.Cache
}

func NewService(db db.DB, cache cache.Cache) Service {
	return &service{
		db:    db,
		cache: cache,
	}
}

func (s *service) Create(req *model.CreateReq) (err error) {
	return errors.WithStack(s.db.Create(req))
}

func (s *service) Update(id uint, req *model.UpdateReq) (err error) {
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

func (s *service) GetRandom20() ([]model.MaximDTO, error) {
	var res []model.MaximDTO
	var err error
	if res, err = s.cache.GetRandom20(); err != nil {
		if errors.Is(err, redis.Nil) {
			if res, err = s.db.GetRandom20(); err != nil {
				return nil, errors.WithStack(err)
			}
			if err = s.cache.SaveRandom20(res); err != nil {
				return nil, errors.WithStack(err)
			}
			return res, err
		}
		return nil, errors.WithStack(err)
	}

	return res, nil
}
