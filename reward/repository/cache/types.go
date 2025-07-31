package cache

import (
	"context"
	"webook/reward/domain"
)

type RewardCache interface {
	GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	SetCachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error
}
