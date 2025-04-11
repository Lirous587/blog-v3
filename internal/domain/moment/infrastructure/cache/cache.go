package cache

import (
	"blog/internal/domain/moment/model"
	"blog/utils"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	SaveRandom20(list []model.MomentDTO) error
	GetRandom20() ([]model.MomentDTO, error)
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &cache{client: client}
}

const (
	momentRandom20Key         = "moment:random20"
	momentRandom20KeyDuration = 2 * time.Hour
)

func (ch *cache) SaveRandom20(list []model.MomentDTO) error {
	ctx := context.Background()
	key := utils.GetRedisKey(momentRandom20Key)
	pipeline := ch.client.Pipeline()

	pipeline.JSONSet(ctx, key, ".", list)
	pipeline.Expire(ctx, key, momentRandom20KeyDuration)

	if _, err := pipeline.Exec(ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ch *cache) GetRandom20() ([]model.MomentDTO, error) {
	ctx := context.Background()
	key := utils.GetRedisKey(momentRandom20Key)

	result, err := ch.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == "" {
		return nil, redis.Nil
	}

	var list []model.MomentDTO
	if err = json.Unmarshal([]byte(result), &list); err != nil {
		return nil, errors.WithStack(err)
	}

	return list, err
}
