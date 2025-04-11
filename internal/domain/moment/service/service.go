package service

import (
	"blog/internal/domain/moment/infrastructure/cache"
	"blog/internal/domain/moment/infrastructure/db"
	"blog/internal/domain/moment/model"
	"blog/pkg/response"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Create(req *model.CreateReq) *response.AppError
	Update(id uint, req *model.UpdateReq) *response.AppError
	Delete(id uint) *response.AppError
	List(req *model.ListReq) (*model.ListRes, *response.AppError)
	GetRandom20() ([]model.MomentDTO, *response.AppError)
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
	if err := s.db.Create(req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) checkRecordIFExistByID(id uint) (*model.Moment, *response.AppError) {
	target, err := s.db.FindByID(id)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}
	if target == nil {
		return nil, response.NewAppError(response.CodeMomentNotFound, err)
	}
	return target, nil
}

func (s *service) Update(id uint, req *model.UpdateReq) *response.AppError {
	if _, err := s.checkRecordIFExistByID(id); err != nil {
		return err
	}

	if err := s.db.Update(id, req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) Delete(id uint) *response.AppError {
	if _, err := s.checkRecordIFExistByID(id); err != nil {
		return err
	}
	if err := s.db.Delete(id); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}

func (s *service) List(req *model.ListReq) (*model.ListRes, *response.AppError) {
	res, err := s.db.List(req)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}
	return res, nil
}

func (s *service) GetRandom20() ([]model.MomentDTO, *response.AppError) {
	var res []model.MomentDTO
	var err error
	if res, err = s.cache.GetRandom20(); err != nil {
		if errors.Is(err, redis.Nil) {
			if res, err = s.db.GetRandom20(); err != nil {
				return nil, response.NewAppError(response.CodeServerError, err)
			}
			if err = s.cache.SaveRandom20(res); err != nil {
				return nil, response.NewAppError(response.CodeServerError, err)
			}
			return res, nil
		}
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	return res, nil
}
