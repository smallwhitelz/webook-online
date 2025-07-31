package cache

import (
	"context"
	"webook/account/domain"
)

type AccountCache interface {
	GetUnique(ctx context.Context, cr domain.Credit) error
	SetUnique(ctx context.Context, cr domain.Credit) error
}
