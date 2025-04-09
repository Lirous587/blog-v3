package service

import (
	"blog/internal/domain/friendLink/model"
	"blog/internal/domain/friendLink/repository/cache"
	"blog/internal/domain/friendLink/repository/db"
	"blog/pkg/response"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Create(req *model.CreateReq) *response.AppError
	Update(id uint, req *model.UpdateReq) *response.AppError
	UpdateStatus(id uint, req *model.UpdateStatusReq) error
	Delete(id uint) error
	List(req *model.ListReq) (res *model.ListRes, err error)
	GetPublicRandom20() ([]model.MaximDTO, error)
	Apply(req *model.ApplyReq) *response.AppError
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

func (s *service) Create(req *model.CreateReq) *response.AppError {
	entity, err := s.db.FindByURL(req.Url)
	if err != nil {
		return response.NewAppError(response.CodeDatabaseError, err)
	}
	if entity != nil {
		return response.NewAppError(response.CodeFriendLinkUrlDuplicate, err)
	}

	if err := s.db.Create(req); err != nil {
		return response.NewAppError(response.CodeDatabaseError, err)
	}
	return nil
}

func (s *service) Update(id uint, req *model.UpdateReq) *response.AppError {
	entity, err := s.db.FindByURL(req.Url)
	if err != nil {
		return response.NewAppError(response.CodeDatabaseError, err)
	}
	if entity != nil {
		return response.NewAppError(response.CodeFriendLinkUrlDuplicate, err)
	}

	if err := s.db.Update(id, req); err != nil {
		return response.NewAppError(response.CodeDatabaseError, err)
	}
	return nil
}

func (s *service) UpdateStatus(id uint, req *model.UpdateStatusReq) error {
	return errors.WithStack(s.db.UpdateStatus(id, req))
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

func (s *service) GetPublicRandom20() ([]model.MaximDTO, error) {
	var res []model.MaximDTO
	var err error
	if res, err = s.cache.GetPublicRandom20(); err != nil {
		if errors.Is(err, redis.Nil) {
			if res, err = s.db.GetPublicRandom20(); err != nil {
				return nil, errors.WithStack(err)
			}
			if err = s.cache.SavePublicRandom20(res); err != nil {
				return nil, errors.WithStack(err)
			}
			return res, err
		}
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (s *service) Apply(req *model.ApplyReq) *response.AppError {
	//TODO implement me
	panic("implement me")
}
