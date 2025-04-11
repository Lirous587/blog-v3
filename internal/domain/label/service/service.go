package service

import (
	"blog/internal/domain/label/model"
	"blog/internal/domain/label/repository/cache"
	"blog/internal/domain/label/repository/db"
	"blog/pkg/response"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Create(req *model.CreateReq) *response.AppError
	Update(id uint, req *model.UpdateReq) *response.AppError
	Delete(id uint) *response.AppError
	List(req *model.ListReq) (*model.ListRes, *response.AppError)
	FindByIDs(ids []uint) ([]model.Label, *response.AppError)
	GetAllWithEssayCount() ([]model.LabelDTO, *response.AppError)
	GetAll() ([]model.LabelDTO, *response.AppError)
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
	target, err := s.db.FindByName(req.Name)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	if target != nil && target.Name == req.Name {
		return response.NewAppError(response.CodeLabelNameDuplicate, errors.New("标签名重复"))
	}

	if err = s.db.Create(req); err != nil {
		return response.NewAppError(response.CodeServerError, errors.New("标签名重复"))
	}

	return nil
}

func (s *service) Update(id uint, req *model.UpdateReq) *response.AppError {
	target, err := s.db.FindByName(req.Name)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	if target != nil && target.ID != id && req.Name == target.Name {
		return response.NewAppError(response.CodeLabelNameDuplicate, errors.New("存在重复的label记录"))
	}

	if err = s.db.Update(id, req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) Delete(id uint) *response.AppError {
	target, err := s.db.FindByID(id)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	if target == nil {
		return response.NewAppError(response.CodeLabelNotFound, nil)
	}

	if err = s.db.Delete(id); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}

func (s *service) List(req *model.ListReq) (*model.ListRes, *response.AppError) {
	res, err := s.db.List(req)
	return res, response.NewAppError(response.CodeServerError, err)
}

func (s *service) FindByIDs(ids []uint) ([]model.Label, *response.AppError) {
	labels, err := s.db.FindByIDs(ids)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}
	return labels, nil
}

func (s *service) GetAllWithEssayCount() ([]model.LabelDTO, *response.AppError) {
	var res []model.LabelDTO
	var err error
	if res, err = s.cache.GetAllWithEssayCount(); err != nil {
		if errors.Is(err, redis.Nil) {
			if res, err = s.db.GetAllWithEssayCount(); err != nil {
				return nil, response.NewAppError(response.CodeServerError, err)

			}
			if err = s.cache.SaveAllWithEssayCount(res); err != nil {
				return nil, response.NewAppError(response.CodeServerError, err)
			}
			return res, nil
		}
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	return res, nil
}

func (s *service) GetAll() ([]model.LabelDTO, *response.AppError) {
	res, err := s.db.GetAll()

	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	return res, nil
}
