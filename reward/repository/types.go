package repository

import (
	"context"
	"webook/reward/domain"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, r domain.Reward) (int64, error)
	GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	SetCachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error
	UpdateReward(ctx context.Context, rid int64, status domain.RewardStatus) error
	GetReward(ctx context.Context, rid int64) (domain.Reward, error)
}
