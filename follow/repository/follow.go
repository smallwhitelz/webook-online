package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"webook/follow/domain"
	"webook/follow/repository/cache"
	"webook/follow/repository/dao"
	"webook/pkg/logger"
)

type FollowRepository interface {
	AddFollowRelation(ctx context.Context, followRelation domain.FollowRelation) error
	CancelFollow(ctx context.Context, followRelation domain.FollowRelation) error
	GetFollowee(ctx context.Context, follower int64, offset int64, limit int64) ([]domain.FollowRelation, error)
	FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error)
	GetFollowStatic(ctx context.Context, uid int64) (domain.FollowStatics, error)
	GetFollower(ctx context.Context, followee int64, offset int64, limit int64) ([]domain.FollowRelation, error)
}

type FollowRelationRepository struct {
	dao   dao.FollowDAO
	cache cache.FollowCache
	l     logger.LoggerV1
}

func NewFollowRelationRepository(dao dao.FollowDAO, cache cache.FollowCache, l logger.LoggerV1) FollowRepository {
	return &FollowRelationRepository{dao: dao, cache: cache, l: l}
}

// GetFollower 获取某个人粉丝列表
func (f *FollowRelationRepository) GetFollower(ctx context.Context, followee int64, offset int64, limit int64) ([]domain.FollowRelation, error) {
	fList, err := f.dao.GetFollowerList(ctx, followee, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(fList, func(idx int, src dao.FollowRelation) domain.FollowRelation {
		return f.toDomain(src)
	}), nil
}

// GetFollowStatic 获取个人关注了多少人，以及粉丝的数量
func (f *FollowRelationRepository) GetFollowStatic(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	// 快路径
	res, err := f.cache.StaticInfo(ctx, uid)
	if err == nil {
		return res, nil
	}
	// 慢路径
	// 这里也可以考虑引入异步机制，比如errgroup
	// 没有就去数据库查
	res.Followers, err = f.dao.CntFollower(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}
	res.Followees, err = f.dao.CntFollowee(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}
	// 回写到缓存里
	err = f.cache.SetStaticInfo(ctx, uid, res)
	if err != nil {
		f.l.Error("回写用户关注人数量到缓存失败", logger.Error(err), logger.Int64("uid", uid))
	}
	return res, nil
}

func (f *FollowRelationRepository) FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error) {
	fr, err := f.dao.FollowRelationDetail(ctx, follower, followee)
	if err != nil {
		return domain.FollowRelation{}, err
	}
	return f.toDomain(fr), nil
}

func (f *FollowRelationRepository) GetFollowee(ctx context.Context, follower int64, offset int64, limit int64) ([]domain.FollowRelation, error) {
	// 如果想要缓存，可以在这里缓存第一页
	fList, err := f.dao.GetFolloweeList(ctx, follower, offset, limit)
	if err != nil {
		return nil, err
	}
	var res []domain.FollowRelation
	for _, fr := range fList {
		res = append(res, f.toDomain(fr))
	}
	return res, nil
}

func (f *FollowRelationRepository) AddFollowRelation(ctx context.Context, followRelation domain.FollowRelation) error {
	err := f.dao.CreateFollowRelation(ctx, f.toEntity(followRelation))
	if err != nil {
		return err
	}
	// 更新缓存里面的关注了多少人以及有多少粉丝的计数
	// 这里已经改为从过canal->kafka进行修改缓存中的计数
	return nil
}

func (f *FollowRelationRepository) CancelFollow(ctx context.Context, followRelation domain.FollowRelation) error {
	err := f.dao.UpdateStatus(ctx, f.toEntity(followRelation), dao.FollowRelationStatusInactive)
	if err != nil {
		return err
	}
	// 这里已经改为从过canal->kafka进行修改缓存中的计数
	return nil
}

func (f *FollowRelationRepository) toEntity(followRelation domain.FollowRelation) dao.FollowRelation {
	return dao.FollowRelation{
		Follower: followRelation.Follower,
		Followee: followRelation.Followee,
	}
}

func (f *FollowRelationRepository) toDomain(fr dao.FollowRelation) domain.FollowRelation {
	return domain.FollowRelation{
		Followee: fr.Followee,
		Follower: fr.Follower,
	}
}
