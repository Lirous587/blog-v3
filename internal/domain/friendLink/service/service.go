package service

import (
	"blog/internal/domain/friendLink/infrastructure/cache"
	"blog/internal/domain/friendLink/infrastructure/db"
	"blog/internal/domain/friendLink/infrastructure/notifier"
	"blog/internal/domain/friendLink/model"
	"blog/pkg/response"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Create(req *model.CreateReq) *response.AppError
	Update(id uint, req *model.UpdateReq) *response.AppError
	Delete(id uint, req *model.DeleteReq) *response.AppError
	List(req *model.ListReq) (res *model.ListRes, err error)
	GetPublishedRandom20() ([]model.FriendLinkDTO, error)
	Apply(req *model.ApplyReq) *response.AppError
	Approve(id uint) *response.AppError
	Reject(id uint) *response.AppError
}

type service struct {
	db       db.DB
	cache    cache.Cache
	notifier notifier.Notifier
}

func NewService(db db.DB, cache cache.Cache, notifier notifier.Notifier) Service {
	return &service{
		db:       db,
		cache:    cache,
		notifier: notifier,
	}
}

func (s *service) Create(req *model.CreateReq) *response.AppError {
	entity, err := s.db.FindByURL(req.Url)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	if entity != nil {
		return response.NewAppError(response.CodeFriendLinkUrlDuplicate, err)
	}

	if err = s.db.Create(req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}

func (s *service) checkRecordIFExistByID(id uint) (*model.FriendLink, *response.AppError) {
	target, err := s.db.FindByID(id)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}
	if target == nil {
		return nil, response.NewAppError(response.CodeFriendLinkNotFound, err)
	}
	return target, nil
}

func (s *service) Update(id uint, req *model.UpdateReq) *response.AppError {
	// 检查要更新的记录是否存在
	if _, appErr := s.checkRecordIFExistByID(id); appErr != nil {
		return appErr
	}
	// 检查URL唯一性
	entity, err := s.db.FindByURL(req.Url)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	if entity != nil && entity.ID != id {
		return response.NewAppError(response.CodeFriendLinkUrlDuplicate, errors.New("该记录已存在"))
	}

	if err = s.db.Update(id, req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) Delete(id uint, req *model.DeleteReq) *response.AppError {
	// 检查要更新的记录是否存在
	target, appErr := s.checkRecordIFExistByID(id)
	if appErr != nil {
		return appErr
	}

	if err := s.notifier.SendDeleteNotification(target, req.Reason); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	if err := s.db.Delete(id); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) List(req *model.ListReq) (res *model.ListRes, err error) {
	res, err = s.db.List(req)
	return res, errors.WithStack(err)
}

func (s *service) GetPublishedRandom20() ([]model.FriendLinkDTO, error) {
	var res []model.FriendLinkDTO
	var err error
	if res, err = s.cache.GetPublishedRandom20(); err != nil {
		if errors.Is(err, redis.Nil) {
			if res, err = s.db.GetPublishedRandom20(); err != nil {
				return nil, errors.WithStack(err)
			}
			if err = s.cache.SavePublishedRandom20(res); err != nil {
				return nil, errors.WithStack(err)
			}
			return res, err
		}
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (s *service) Apply(req *model.ApplyReq) *response.AppError {
	entity, err := s.db.FindByURL(req.Url)
	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	if entity != nil {
		return response.NewAppError(response.CodeFriendLinkUrlDuplicate, err)
	}

	if err = s.db.Apply(req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) Approve(id uint) *response.AppError {
	// 检查记录是否存在
	target, appErr := s.checkRecordIFExistByID(id)
	if appErr != nil {
		return appErr
	}

	if target.Status == model.StatusPublished {
		return response.NewAppError(response.CodeIllegalOperation, errors.New("友链状态为published,无法进行审核操作"))
	}

	// 发邮件
	if err := s.notifier.SendApprovalNotification(target); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	// 修改db
	req := &model.UpdateStatusReq{Status: model.StatusPublished}
	if err := s.db.UpdateStatus(id, req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
}

func (s *service) Reject(id uint) *response.AppError {
	// 检查记录是否存在
	target, appErr := s.checkRecordIFExistByID(id)
	if appErr != nil {
		return appErr
	}

	if target.Status == model.StatusPublished {
		return response.NewAppError(response.CodeIllegalOperation, errors.New("友链状态为published,无法进行审核操作"))
	}

	// 发邮件
	if err := s.notifier.SendRejectionNotification(target); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	// 删除申请记录
	if err := s.db.Delete(id); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}
