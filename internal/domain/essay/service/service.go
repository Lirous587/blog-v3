package service

import (
	"blog/internal/domain/essay/infrastructure/cache"
	"blog/internal/domain/essay/infrastructure/db"
	"blog/internal/domain/essay/model"
	model2 "blog/internal/domain/label/model"
	labelService "blog/internal/domain/label/service"
	"blog/pkg/response"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	Create(req *model.CreateReq) *response.AppError
	Update(id uint, req *model.UpdateReq) *response.AppError
	Delete(id uint) *response.AppError
	Get(id uint) (*model.EssayDTO, *response.AppError)
	GetNears(id uint) ([]model.EssayDTO, *response.AppError)
	List(req *model.ListReq) (*model.ListRes, *response.AppError)
	GetTimelines() ([]model.Timeline, *response.AppError)
}

type service struct {
	db           db.DB
	cache        cache.Cache
	labelService labelService.Service
}

func NewService(db db.DB, cache cache.Cache, labelService labelService.Service) Service {
	return &service{db: db, cache: cache, labelService: labelService}
}

func (s *service) findLabelsByIDs(ids []uint) ([]model2.Label, *response.AppError) {
	labels, err := s.labelService.FindByIDs(ids)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	if len(labels) != len(ids) {
		return nil, response.NewAppError(response.CodeServerError, errors.New("无效的输入标签"))
	}
	return labels, nil
}

func (s *service) Create(req *model.CreateReq) *response.AppError {
	labels, err := s.findLabelsByIDs(req.LabelIDs)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	if err := s.db.Create(req, labels); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) Update(id uint, req *model.UpdateReq) *response.AppError {
	labels, err := s.findLabelsByIDs(req.LabelIDs)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	if err := s.db.Update(id, req, labels); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}

func (s *service) Delete(id uint) *response.AppError {
	_, err := s.db.FindByID(id)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	if err := s.db.Delete(id); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}

func (s *service) Get(id uint) (*model.EssayDTO, *response.AppError) {
	errGroup := errgroup.Group{}

	var res *model.EssayDTO
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
		res, err = s.db.GetByID(id)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	res.VisitedTimes = vt

	return res, nil
}

func (s *service) GetNears(id uint) ([]model.EssayDTO, *response.AppError) {
	res, err := s.db.GetNearsByID(id)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}
	return res, nil
}

func (s *service) List(req *model.ListReq) (*model.ListRes, *response.AppError) {
	res, err := s.db.List(req)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	ids := make([]uint, len(res.List))
	for i := range ids {
		ids[i] = res.List[i].ID
	}

	vtMap, err := s.cache.GetNVisitedTimes(ids)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
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

func (s *service) GetTimelines() ([]model.Timeline, *response.AppError) {
	res, err := s.cache.GetTimeline()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			res, err = s.db.GetTimelines()
			if err != nil {
				return nil, response.NewAppError(response.CodeServerError, err)
			}

			if err = s.cache.SaveTimeline(res); err != nil {
				return nil, response.NewAppError(response.CodeServerError, err)
			}
			return res, nil
		}
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	return res, nil
}
