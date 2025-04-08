package cache

import (
	"blog/internal/domain/essay/model"
	"blog/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Cache interface {
	SaveVisitedTimes(id, vt uint) error
	GetVisitedTimes(id uint) (uint, error)
	GetNVisitedTimes(ids []uint) (map[uint]uint, error)
	GetAllVt() (map[uint]uint, error)
	SaveTimeline(data []model.Timeline) error
	GetTimeline() ([]model.Timeline, error)
}

type cache struct {
	client *redis.Client
}

const (
	essayVisitedTimesKey     = "essay:visitedTimesMap"
	essayTimelineKey         = "essay:timeline"
	essayTimelineKeyDuration = 2 * time.Hour
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

func (ch *cache) GetNVisitedTimes(ids []uint) (map[uint]uint, error) {
	key := utils.GetRedisKey(essayVisitedTimesKey)

	idsStr := make([]string, 0, len(ids))

	for _, id := range ids {
		idsStr = append(idsStr, fmt.Sprintf("%d", id))
	}

	vts, err := ch.client.HMGet(context.Background(), key, idsStr...).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	vtMap := make(map[uint]uint, len(vts))

	for i, vt := range vts {
		if vt != nil {
			vtUint, _ := strconv.ParseUint(vt.(string), 10, 64)
			vtMap[ids[i]] = uint(vtUint)
		} else {
			vtMap[ids[i]] = 0
		}
	}

	return vtMap, nil
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

func (ch *cache) SaveTimeline(list []model.Timeline) error {
	key := utils.GetRedisKey(essayTimelineKey)
	ctx := context.Background()
	pipeline := ch.client.Pipeline()

	pipeline.JSONSet(ctx, key, ".", list)

	pipeline.Expire(ctx, key, essayTimelineKeyDuration)

	if _, err := pipeline.Exec(ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ch *cache) GetTimeline() ([]model.Timeline, error) {
	key := utils.GetRedisKey(essayTimelineKey)
	ctx := context.Background()

	result, err := ch.client.JSONGet(ctx, key).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if result == "" {
		return nil, redis.Nil
	}

	var res []model.Timeline

	if err = json.Unmarshal([]byte(result), &res); err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}
