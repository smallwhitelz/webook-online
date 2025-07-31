package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/account/domain"
)

type AccountRedisCache struct {
	client redis.Cmdable
}

func NewAccountRedisCache(client redis.Cmdable) AccountCache {
	return &AccountRedisCache{client: client}
}

func (a *AccountRedisCache) GetUnique(ctx context.Context, cr domain.Credit) error {
	res, err := a.client.Exists(ctx, a.key(cr.Biz, cr.BizId)).Result()
	if err != nil {
		return err
	}
	if res > 0 {
		return errors.New("该业务已经处理过了")
	}
	return nil
}

func (a *AccountRedisCache) SetUnique(ctx context.Context, cr domain.Credit) error {
	return a.client.Set(ctx, a.key(cr.Biz, cr.BizId), "", time.Minute*30).Err()
}

func (a *AccountRedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("credit:biz:%s_%d", biz, bizId)
}
