package service

import (
	"context"
	"webook/follow/domain"
	"webook/follow/repository"
)

type FollowService interface {
	Follow(ctx context.Context, follower, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
	// GetFollowee 获得某个人的关注列表
	GetFollowee(ctx context.Context, follower int64, offset, limit int64) ([]domain.FollowRelation, error)
	// FollowInfo 点进关注的人的文章或者视频详情里，展示关注的状态
	FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error)

	// GetFollower 获取某个人的粉丝列表
	GetFollower(ctx context.Context, followee int64, offset, limit int64) ([]domain.FollowRelation, error)
	// GetFollowStatic 获取个人关注了多少人，以及粉丝的数量
	GetFollowStatic(ctx context.Context, uid int64) (domain.FollowStatics, error)
}

type FollowRelationService struct {
	repo repository.FollowRepository
}

func NewFollowRelationService(repo repository.FollowRepository) FollowService {
	return &FollowRelationService{repo: repo}
}

func (f *FollowRelationService) GetFollower(ctx context.Context, followee int64, offset, limit int64) ([]domain.FollowRelation, error) {
	return f.repo.GetFollower(ctx, followee, offset, limit)
}

func (f *FollowRelationService) GetFollowStatic(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	return f.GetFollowStatic(ctx, uid)
}

func (f *FollowRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	return f.repo.FollowInfo(ctx, follower, followee)
}

func (f *FollowRelationService) GetFollowee(ctx context.Context, follower int64, offset, limit int64) ([]domain.FollowRelation, error) {
	return f.repo.GetFollowee(ctx, follower, offset, limit)
}

func (f *FollowRelationService) Follow(ctx context.Context, follower, followee int64) error {
	return f.repo.AddFollowRelation(ctx, domain.FollowRelation{
		Followee: followee,
		Follower: follower,
	})
}

func (f *FollowRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	return f.repo.CancelFollow(ctx, domain.FollowRelation{
		Followee: followee,
		Follower: follower,
	})
}
