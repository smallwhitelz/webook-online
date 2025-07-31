package repository

import (
	"context"
	"webook/reward/domain"
	"webook/reward/repository/cache"
	"webook/reward/repository/dao"
)

type rewardRepository struct {
	dao   dao.RewardDAO
	cache cache.RewardCache
}

func NewRewardRepository(dao dao.RewardDAO, cache cache.RewardCache) RewardRepository {
	return &rewardRepository{dao: dao, cache: cache}
}

func (repo *rewardRepository) GetReward(ctx context.Context, rid int64) (domain.Reward, error) {
	rw, err := repo.dao.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	return repo.toDomain(rw), nil
}

func (repo *rewardRepository) UpdateReward(ctx context.Context, rid int64, status domain.RewardStatus) error {
	return repo.dao.UpdateReward(ctx, rid, status)
}

func (repo *rewardRepository) CreateReward(ctx context.Context, r domain.Reward) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(r))
}

func (repo *rewardRepository) GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	return repo.cache.GetCachedCodeURL(ctx, r)
}

func (repo *rewardRepository) SetCachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error {
	return repo.cache.SetCachedCodeURL(ctx, cu, r)
}

func (repo *rewardRepository) toEntity(r domain.Reward) dao.Reward {
	return dao.Reward{
		Biz:       r.Biz,
		BizId:     r.BizId,
		BizName:   r.BizName,
		TargetUid: r.TargetUid,
		Status:    r.Status.AsUint8(),
		Uid:       r.Uid,
		Amount:    r.Amt,
	}
}

func (repo *rewardRepository) toDomain(rw dao.Reward) domain.Reward {
	return domain.Reward{
		Id:        rw.Id,
		Biz:       rw.Biz,
		BizId:     rw.BizId,
		BizName:   rw.BizName,
		TargetUid: rw.TargetUid,
		Uid:       rw.Uid,
		Amt:       rw.Amount,
		Status:    domain.RewardStatus(rw.Status),
	}
}
