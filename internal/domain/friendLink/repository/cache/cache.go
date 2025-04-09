package cache

import (
	"blog/internal/domain/friendLink/model"
	"blog/utils"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	SavePublicRandom20(list []model.MaximDTO) error
	GetPublicRandom20() ([]model.MaximDTO, error)
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &cache{client: client}
}

const (
	maximRandom20Key         = "friendLink:random20"
	maximRandom20KeyDuration = 2 * time.Hour
)

func (ch *cache) SavePublicRandom20(list []model.MaximDTO) error {
	ctx := context.Background()
	key := utils.GetRedisKey(maximRandom20Key)
	pipeline := ch.client.Pipeline()

	pipeline.JSONSet(ctx, key, ".", list)
	pipeline.Expire(ctx, key, maximRandom20KeyDuration)

	if _, err := pipeline.Exec(ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ch *cache) GetPublicRandom20() ([]model.MaximDTO, error) {
	ctx := context.Background()
	key := utils.GetRedisKey(maximRandom20Key)

	result, err := ch.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == "" {
		return nil, redis.Nil
	}

	var list []model.MaximDTO
	if err = json.Unmarshal([]byte(result), &list); err != nil {
		return nil, errors.WithStack(err)
	}

	return list, err
}
