package cache

import (
	"blog/internal/domain/essay/model"
	"blog/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Cache interface {
	GetVisitedTimes(id uint) (uint, error)
	SaveVisitedTimes(id, vt uint) error

	GetAllVt() (map[uint]uint, error)
	GetDates() ([]string, error)
}

type cache struct {
	client *redis.Client
}

const (
	essayVisitedTimesKey = "essay:visitedTimesMap"
	essayDatesKey        = "essay:dates"
)

func NewCache(client *redis.Client) Cache {
	return &cache{client: client}
}

func (ch *cache) GetVisitedTimes(id uint) (uint, error) {
	idStr := fmt.Sprintf("%d", id)
	key := utils.GetRedisKey(essayVisitedTimesKey)

	vt, err := ch.client.HIncrBy(context.Background(), key, idStr, 1).Result()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return uint(vt), nil
}

func (ch *cache) SaveVisitedTimes(id, vt uint) error {
	idStr := fmt.Sprintf("%d", id)
	key := utils.GetRedisKey(essayVisitedTimesKey)

	if err := ch.client.HSet(context.Background(), key, idStr, vt).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ch *cache) GetAllVt() (map[uint]uint, error) {
	key := utils.GetRedisKey(essayVisitedTimesKey)
	ctx := context.Background()

	result, err := ch.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	vtMap := make(map[uint]uint, len(result))
	for idStr, vtStr := range result {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		vt, _ := strconv.ParseUint(vtStr, 10, 64)
		vtMap[uint(id)] = uint(vt)
	}

	return vtMap, nil
}

func (ch *cache) GetDates() ([]string, error) {
	key := utils.GetRedisKey(essayDatesKey)

	res, err := ch.client.LRange(context.Background(), key, 0, -1).Result()

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
			pipe := ch.client.Pipeline()
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
