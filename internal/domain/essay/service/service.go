package service

import (
	"blog/internal/domain/essay/model"
	"blog/internal/domain/essay/repository/cache"
	"blog/internal/domain/essay/repository/db"
	model2 "blog/internal/domain/label/model"
	labelService "blog/internal/domain/label/service"
	"github.com/pkg/errors"
)

type Service interface {
	Create(req *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	Get(id uint) (*model.GetRes, error)
	List(req *model.ListReq) (*model.ListRes, error)
	GetTimelines(req *model.TimelineReq) (*model.TimelineRes, error)
}

type service struct {
	db           db.DB
	cache        cache.Cache
	labelService labelService.Service
}

func NewService(db db.DB, cache cache.Cache, labelService labelService.Service) Service {
	return &service{db: db, cache: cache, labelService: labelService}
}

func (s *service) findLabelsByIDs(ids []uint) ([]model2.Label, error) {
	labels, err := s.labelService.FindByIDs(ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(labels) != len(ids) {
		return nil, errors.New("无效的输入标签")
	}
	return labels, nil
}

func (s *service) Create(req *model.CreateReq) error {
	labels, err := s.findLabelsByIDs(req.LabelIDs)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(s.db.Create(req, labels))
}

func (s *service) Update(id uint, req *model.UpdateReq) (err error) {
	labels, err := s.findLabelsByIDs(req.LabelIDs)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(s.db.Update(id, req, labels))
}

func (s *service) Delete(id uint) (err error) {
	_, err = s.db.FindByID(id)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(s.db.Delete(id))
}

func (s *service) Get(id uint) (*model.GetRes, error) {
	res, err := s.db.Get(id)
	return res, errors.WithStack(err)
}

func (s *service) List(req *model.ListReq) (*model.ListRes, error) {
	res, err := s.db.List(req)
	return res, errors.WithStack(err)
}

func (s *service) GetTimelines(req *model.TimelineReq) (*model.TimelineRes, error) {
	//dates, err := s.cache.GetDates()
	//if err != nil {
	//	return nil, errors.WithStack(err)
	//}

	//res, err := s.repo.GetTimelines(req)
	//return res, errors.WithStack(err)
	return nil, nil
}
