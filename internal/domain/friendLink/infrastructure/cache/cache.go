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
	SavePublishedRandom20(list []model.FriendLinkDTO) error
	GetPublishedRandom20() ([]model.FriendLinkDTO, error)
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &cache{client: client}
}

const (
	random20Key         = "friendLink:random20"
	random20KeyDuration = 2 * time.Hour
)

func (ch *cache) SavePublishedRandom20(list []model.FriendLinkDTO) error {
	ctx := context.Background()
	key := utils.GetRedisKey(random20Key)
	pipeline := ch.client.Pipeline()

	pipeline.JSONSet(ctx, key, ".", list)
	pipeline.Expire(ctx, key, random20KeyDuration)

	if _, err := pipeline.Exec(ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ch *cache) GetPublishedRandom20() ([]model.FriendLinkDTO, error) {
	ctx := context.Background()
	key := utils.GetRedisKey(random20Key)

	result, err := ch.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == "" {
		return nil, redis.Nil
	}

	var list []model.FriendLinkDTO
	if err = json.Unmarshal([]byte(result), &list); err != nil {
		return nil, errors.WithStack(err)
	}

	return list, err
}
