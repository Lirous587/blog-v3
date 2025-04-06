package db

import (
	"blog/internal/domain/essay/model"
	labelModel "blog/internal/domain/label/model"
	"blog/utils"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type DB interface {
	FindByID(id uint) (*model.Essay, error)
	Create(req *model.CreateReq, labels []labelModel.Label) error
	Update(id uint, req *model.UpdateReq, labels []labelModel.Label) error
	Delete(id uint) error
	Get(id uint) (*model.GetRes, error)
	List(res *model.ListReq) (*model.ListRes, error)
	FindVTsByIDs(ids []uint) (map[uint]uint, error)
	SaveVTsByIDs(idVtMap map[uint]uint) error
	GetTimelines(req *model.TimelineReq) (*model.TimelineRes, error)
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
	label := new(model.Essay)
	if err := db.orm.Where("id = ?", id).First(label).Error; err != nil {
		return nil, err
	}
	return label, nil
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

func (db *db) Get(id uint) (*model.GetRes, error) {
	var essay model.Essay
	var previous []model.Essay
	var next []model.Essay

	var errGroup errgroup.Group

	// id对应的文章
	errGroup.Go(func() error {
		return db.orm.Preload("Labels").First(&essay, id).Error
	})

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

	previousDTOs := make([]model.EssayDTO, len(previous))
	nextDTOs := make([]model.EssayDTO, len(next))

	for i := range previous {
		previousDTOs[i] = *previous[i].ConvertToDTO()
	}

	for i := range next {
		nextDTOs[i] = *next[i].ConvertToDTO()
	}

	res := model.GetRes{
		EssayDTO: *essay.ConvertToDTO(),
		Previous: previousDTOs,
		Next:     nextDTOs,
	}
	return &res, errGroup.Wait()
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

	// 准备批量更新的数据
	var updates []map[string]interface{}
	for id, vt := range idVtMap {
		updates = append(updates, map[string]interface{}{
			"id":            id,
			"visited_times": vt,
		})
	}

	// 执行批量更新
	return db.orm.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Essay{}).Updates(updates)
		return errors.WithStack(result.Error)
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

func (db *db) GetTimelines(req *model.TimelineReq) (*model.TimelineRes, error) {
	return nil, nil
}
