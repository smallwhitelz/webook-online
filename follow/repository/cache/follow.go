package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"webook/follow/domain"
)

type FollowCache interface {
	StaticInfo(ctx context.Context, uid int64) (domain.FollowStatics, error)
	SetStaticInfo(ctx context.Context, uid int64, followStatics domain.FollowStatics) error
	Follow(ctx context.Context, follower int64, followee int64) error
	CancelFollow(ctx context.Context, follower int64, followee int64) error
}

const (
	// 被多少人关注
	fieldFollowerCnt = "follower_cnt"
	// 关注了多少人
	fieldFolloweeCnt = "followee_cnt"
)

type FollowRedisCache struct {
	client redis.Cmdable
}

func NewFollowRedisCache(client redis.Cmdable) FollowCache {
	return &FollowRedisCache{client: client}
}

func (f *FollowRedisCache) StaticInfo(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	data, err := f.client.HGetAll(ctx, f.staticsKey(uid)).Result()
	if err != nil {
		return domain.FollowStatics{}, err
	}
	// key 不存在，没有缓存的数据
	if len(data) == 0 {
		return domain.FollowStatics{}, redis.Nil
	}
	var res domain.FollowStatics
	res.Followers, _ = strconv.ParseInt(data[fieldFollowerCnt], 10, 64)
	res.Followees, _ = strconv.ParseInt(data[fieldFolloweeCnt], 10, 64)
	return res, nil
}

func (f *FollowRedisCache) SetStaticInfo(ctx context.Context, uid int64, followStatics domain.FollowStatics) error {
	return f.client.HSet(ctx, f.staticsKey(uid), fieldFollowerCnt, followStatics.Followers, fieldFolloweeCnt, followStatics.Followees).Err()
}

func (f *FollowRedisCache) Follow(ctx context.Context, follower int64, followee int64) error {
	return f.updateStaticsInfo(ctx, follower, followee, 1)
}

func (f *FollowRedisCache) CancelFollow(ctx context.Context, follower int64, followee int64) error {
	return f.updateStaticsInfo(ctx, follower, followee, 1)
}

func (f *FollowRedisCache) staticsKey(uid int64) string {
	return fmt.Sprintf("follow:statics:%d", uid)
}

func (f *FollowRedisCache) updateStaticsInfo(ctx context.Context, follower int64, followee int64, delta int64) error {
	// 用到redis的事务
	// 但是没办法达到acid的特性，这里只是保证这两个命令同时过去到达redis，中间不会有人插队，
	// 但是没办法保证两个都执行成功或者都不成功
	// 任何中间件想要达到acid的特性就要看有没有类似mysql的redolog和undolog的特性
	// 显然 redis是没有的
	// 不能A关注B了 结果A的关注人多了，B的粉丝没变
	tx := f.client.TxPipeline()
	// 这两个操作只是记录了一下，还没有发过去
	f.client.HIncrBy(ctx, f.staticsKey(follower), fieldFolloweeCnt, delta)
	f.client.HIncrBy(ctx, f.staticsKey(followee), fieldFollowerCnt, delta)
	// 发过去redis执行，并且返回了结果
	_, err := tx.Exec(ctx)
	return err
}
