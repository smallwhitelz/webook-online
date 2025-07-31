package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"webook/tag/domain"
)

type TagCache interface {
	Append(ctx context.Context, uid int64, tags ...domain.Tag) error
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
}

type TagRedisCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewTagRedisCache(client redis.Cmdable) TagCache {
	return &TagRedisCache{client: client, expiration: time.Minute * 30}
}

func (t *TagRedisCache) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	data, err := t.client.HGetAll(ctx, t.userTagsKey(uid)).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, redis.Nil
	}
	res := make([]domain.Tag, 0, len(data))
	for _, val := range data {
		var tag domain.Tag
		err := json.Unmarshal([]byte(val), &tag)
		if err != nil {
			return nil, err
		}
		res = append(res, tag)
	}
	return res, nil
}

func (t *TagRedisCache) Append(ctx context.Context, uid int64, tags ...domain.Tag) error {
	// 要放我的标签
	// list ,hash,set ,sorted set
	key := t.userTagsKey(uid)
	// 对同一个用户操作多个标签
	pip := t.client.Pipeline()
	for _, tag := range tags {
		val, err := json.Marshal(tag)
		if err != nil {
			return err
		}
		// uid = > tid_0 -> t0
		//         tid_1 -> t1
		pip.HMSet(ctx, key, strconv.FormatInt(tag.Id, 10), val)
	}
	// 第一次进来，可能没有过期时间，所以需要我们设计
	pip.Expire(ctx, key, t.expiration)
	_, err := pip.Exec(ctx)
	return err
}

func (t *TagRedisCache) userTagsKey(uid int64) string {
	return fmt.Sprintf("tag:user_tags:%d", uid)
}
