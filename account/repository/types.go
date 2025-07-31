package repository

import (
	"context"
	"webook/account/domain"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, cr domain.Credit) error
	CheckUnique(ctx context.Context, cr domain.Credit) error
	SetUnique(ctx context.Context, cr domain.Credit) error
}
