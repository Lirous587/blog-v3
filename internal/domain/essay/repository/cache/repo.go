package cache

import (
	"blog/internal/domain/essay/model"
	"blog/utils"
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	//FindByID(id uint) (*model.Essay, error)
	//Create(req *model.CreateReq, labels []labelModel.Label) error
	//Update(id uint, req *model.UpdateReq, labels []labelModel.Label) error
	//Delete(id uint) error
	//Get(id uint) (*model.GetRes, error)
	//List(res *model.ListReq) (*model.ListRes, error)
	//GetTimelines(req *model.TimelineReq) (*model.TimelineRes, error)
	GetDates() ([]string, error)
}

type cache struct {
	client *redis.Client
}

const (
	essayDatesKey = "essay:dates"
)

func NewCache(client *redis.Client) Cache {
	return &cache{client: client}
}

func (r *cache) GetDates() ([]string, error) {
	key := utils.GetRedisKey(essayDatesKey)

	res, err := r.client.LRange(context.Background(), key, 0, -1).Result()

	if errors.Is(err, redis.Nil) || len(res) == 0 {
		var essays []model.Essay
		//if err := r.db.Model(&model.Essay{}).Select("created_at").Find(&essays).Error; err != nil {
		//	return nil, errors.WithStack(err)
		//}
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

		// 如果有日期，存入Redis
		if len(uniqueDates) > 0 {
			pipe := r.client.Pipeline()
			// 先删除旧key
			pipe.Del(context.Background(), key)
			// 正确使用RPUSH添加多个元素
			pipe.RPush(context.Background(), key, uniqueDates)

			if _, err = pipe.Exec(context.Background()); err != nil {
				return uniqueDates, errors.WithStack(err) // 返回数据但同时报告Redis错误
			}
			return uniqueDates, nil
		}
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}
