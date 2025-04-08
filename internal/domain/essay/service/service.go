package service

import (
	"blog/internal/domain/essay/model"
	"blog/internal/domain/essay/repository/cache"
	"blog/internal/domain/essay/repository/db"
	model2 "blog/internal/domain/label/model"
	labelService "blog/internal/domain/label/service"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	Create(req *model.CreateReq) error
	Update(id uint, req *model.UpdateReq) error
	Delete(id uint) error
	Get(id uint) (*model.GetRes, error)
	List(req *model.ListReq) (*model.ListRes, error)
	GetTimelines() ([]model.Timeline, error)
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
	errGroup := errgroup.Group{}

	var res *model.GetRes
	var vt uint

	// 从redis里面去添加vt并且返回
	errGroup.Go(func() error {
		var err error
		vt, err = s.cache.GetVisitedTimes(id)
		// redis错误或无记录 db兜底
		if err != nil || vt == 0 {
			ids := []uint{id}
			vtMap, err := s.db.FindVTsByIDs(ids)
			if err != nil {
				return errors.WithStack(err)
			}
			if dbVt, exists := vtMap[id]; exists {
				vt = dbVt
				_ = s.cache.SaveVisitedTimes(id, vt+1)
				vt++
			}
		}
		return nil
	})

	// 去db查数据
	errGroup.Go(func() error {
		var err error
		res, err = s.db.Get(id)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	res.VisitedTimes = vt

	return res, nil
}

func (s *service) List(req *model.ListReq) (*model.ListRes, error) {
	res, err := s.db.List(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ids := make([]uint, len(res.List))
	for i := range ids {
		ids[i] = res.List[i].ID
	}

	vtMap, err := s.cache.GetNVisitedTimes(ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range res.List {
		id := res.List[i].ID
		vt := vtMap[id]
		if vt != 0 {
			res.List[i].VisitedTimes = vt
		}
	}

	return res, nil
}

func (s *service) GetTimelines() ([]model.Timeline, error) {
	res, err := s.cache.GetTimeline()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			res, err = s.db.GetTimelines()
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if err = s.cache.SaveTimeline(res); err != nil {
				return nil, errors.WithStack(err)
			}
			return res, nil
		}
		return nil, errors.WithStack(err)
	}

	return res, nil
}
