package cache

import (
	"blog/internal/domain/label/model"
	"blog/utils"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	GetAllWithEssayCount() ([]model.LabelDTO, error)
	SaveAllWithEssayCount(list []model.LabelDTO) error
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &cache{client: client}
}

const (
	labelEssayCountsKey         = "label:essay_counts"
	labelEssayCountsKeyDuration = 24 * time.Hour
)

func (ch *cache) SaveAllWithEssayCount(list []model.LabelDTO) error {
	ctx := context.Background()
	key := utils.GetRedisKey(labelEssayCountsKey)
	pipeline := ch.client.Pipeline()

	pipeline.JSONSet(ctx, key, ".", list)
	pipeline.Expire(ctx, key, labelEssayCountsKeyDuration)

	if _, err := pipeline.Exec(ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ch *cache) GetAllWithEssayCount() ([]model.LabelDTO, error) {
	ctx := context.Background()
	key := utils.GetRedisKey(labelEssayCountsKey)

	result, err := ch.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == "" {
		return nil, redis.Nil
	}

	var list []model.LabelDTO
	if err := json.Unmarshal([]byte(result), &list); err != nil {
		return nil, errors.WithStack(err)
	}

	return list, err
}
