package db

import (
	"blog/internal/domain/essay/model"
	labelModel "blog/internal/domain/label/model"
	"blog/utils"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type DB interface {
	FindByID(id uint) (*model.Essay, error)
	Create(req *model.CreateReq, labels []labelModel.Label) error
	Update(id uint, req *model.UpdateReq, labels []labelModel.Label) error
	Delete(id uint) error
	GetByID(id uint) (*model.EssayDTO, error)
	GetNearsByID(id uint) ([]model.EssayDTO, error)
	List(res *model.ListReq) (*model.ListRes, error)
	FindVTsByIDs(ids []uint) (map[uint]uint, error)
	SaveVTsByIDs(idVtMap map[uint]uint) error
	GetTimelines() ([]model.Timeline, error)
}

type db struct {
	orm *gorm.DB
}

func NewDB(orm *gorm.DB) DB {
	//var essay model.Essay
	//orm.AutoMigrate(&essay)
	return &db{orm: orm}
}

func (db *db) FindByID(id uint) (*model.Essay, error) {
	essay := new(model.Essay)
	if err := db.orm.Where("id = ?", id).First(essay).Error; err != nil {
		return nil, err
	}
	return essay, nil
}

func (db *db) Create(req *model.CreateReq, labels []labelModel.Label) error {
	essay := model.Essay{
		Name:         req.Name,
		Introduction: req.Introduction,
		Content:      req.Content,
		Priority:     req.Priority,
		PreviewTheme: req.PreviewTheme,
		CodeTheme:    req.CodeTheme,
		ImgUrl:       req.ImgUrl,
		Labels:       labels,
	}
	return db.orm.Create(&essay).Error
}

func (db *db) Update(id uint, req *model.UpdateReq, labels []labelModel.Label) error {
	essay := model.Essay{
		Model:        gorm.Model{ID: id}, // 设置ID是关键
		Name:         req.Name,
		Introduction: req.Introduction,
		Content:      req.Content,
		Priority:     req.Priority,
		PreviewTheme: req.PreviewTheme,
		CodeTheme:    req.CodeTheme,
		ImgUrl:       req.ImgUrl,
	}

	return db.orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Labels").Updates(essay).Error; err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(tx.Model(&essay).Association("Labels").Replace(labels))
	})
}

func (db *db) Delete(id uint) error {
	return db.orm.Transaction(func(tx *gorm.DB) error {
		essay := model.Essay{
			Model: gorm.Model{ID: id},
		}
		if err := tx.Model(&essay).Association("Labels").Clear(); err != nil {
			return errors.WithStack(err)
		}
		return tx.Delete(&essay).Error
	})
}

func (db *db) GetByID(id uint) (*model.EssayDTO, error) {
	var essay model.Essay

	if err := db.orm.Preload("Labels").First(&essay, id).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	res := *essay.ConvertToDTO()

	return &res, nil
}

func (db *db) GetNearsByID(id uint) ([]model.EssayDTO, error) {
	var previous []model.Essay
	var next []model.Essay
	var errGroup errgroup.Group

	// 前面的文章
	errGroup.Go(func() error {
		return db.orm.Where("id < ?", id).
			Order("id DESC").
			Limit(3).
			Preload("Labels").
			Omit("content", "VisitedTimes").
			Find(&previous).Error
	})

	// 后面的文章
	errGroup.Go(func() error {
		return db.orm.Where("id > ?", id).
			Order("id ASC").
			Limit(3).
			Preload("Labels").
			Omit("content", "VisitedTimes").
			Find(&next).Error
	})

	if err := errGroup.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	// 合并结果
	result := make([]model.EssayDTO, 0, len(previous)+len(next))

	for i := range previous {
		result = append(result, *previous[i].ConvertToDTO())
	}

	for i := range next {
		result = append(result, *next[i].ConvertToDTO())
	}

	return result, nil
}

func (db *db) FindVTsByIDs(ids []uint) (map[uint]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	type Result struct {
		ID           uint
		VisitedTimes uint
	}

	var results []Result
	if err := db.orm.Model(model.Essay{}).
		Where("id IN ?", ids).
		Select("id", "visited_times").
		Find(&results).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	vtMap := make(map[uint]uint, len(results))

	for i := range results {
		vtMap[results[i].ID] = results[i].VisitedTimes
	}

	return vtMap, nil
}

func (db *db) SaveVTsByIDs(idVtMap map[uint]uint) error {
	if len(idVtMap) == 0 {
		return nil
	}

	// 构建批量更新语句
	type BatchUpdate struct {
		ID           uint
		VisitedTimes uint
	}

	var updates []BatchUpdate
	for id, vt := range idVtMap {
		updates = append(updates, BatchUpdate{ID: id, VisitedTimes: vt})
	}

	return db.orm.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			if err := tx.Model(&model.Essay{}).
				Where("id = ?", update.ID).
				Updates(map[string]interface{}{"visited_times": update.VisitedTimes}).
				Error; err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}

func (db *db) List(req *model.ListReq) (*model.ListRes, error) {
	var count int64
	var essays []model.Essay

	query := db.orm.Model(&model.Essay{})

	if req.LabelID > 0 {
		query = query.Where("id IN (?)",
			db.orm.Table("essay_labels").
				Select("essay_id").
				Where("label_id = ?", req.LabelID))
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	page, err := utils.ComputePages(count, req.PageSize, req.Page)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	offset, err := utils.ComputeOffset(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = query.Preload("Labels").
		Limit(req.PageSize).
		Offset(offset).
		Omit("content").
		Find(&essays).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	dtos := make([]model.EssayDTO, len(essays))

	for i := range dtos {
		dtos[i] = *essays[i].ConvertToDTO()
	}
	return &model.ListRes{
		List:  dtos,
		Pages: page,
	}, nil
}

func (db *db) GetYearDates() ([]string, error) {
	var essays []model.Essay
	if err := db.orm.Model(&model.Essay{}).Select("created_at").Find(&essays).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	// 使用map去重
	dateMap := make(map[string]struct{})
	for i := range essays {
		month := essays[i].CreatedAt.Format("2006-01")
		dateMap[month] = struct{}{}
	}

	// 转换为切片
	uniqueDates := make([]string, 0, len(dateMap))
	for date := range dateMap {
		uniqueDates = append(uniqueDates, date)
	}
	return uniqueDates, nil
}

func (db *db) GetTimelines() ([]model.Timeline, error) {
	// 拿到所有年份
	dates, err := db.GetYearDates()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	timelines := make([]model.Timeline, 0, len(dates))

	for _, date := range dates {
		startTime, err := time.Parse("2006-01", date)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		nextMonth := startTime.AddDate(0, 1, 0)

		var essays []model.Essay
		if err := db.orm.
			Where("created_at >= ? AND created_at < ?", startTime, nextMonth).
			Select([]string{"id", "name", "created_at"}).
			Find(&essays).Error; err != nil {
			return nil, errors.WithStack(err)
		}

		dtos := make([]model.EssayDTO, len(essays))
		for i := range essays {
			dtos[i] = *essays[i].ConvertToDTO()
		}

		if len(dtos) > 0 {
			timelines = append(timelines, model.Timeline{
				Data:    date,
				Records: dtos,
			})
		}
	}

	return timelines, nil
}
