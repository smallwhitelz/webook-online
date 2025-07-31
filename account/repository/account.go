package repository

import (
	"context"
	"time"
	"webook/account/domain"
	"webook/account/repository/cache"
	"webook/account/repository/dao"
)

type accountRepository struct {
	dao   dao.AccountDAO
	cache cache.AccountCache
}

func NewAccountRepository(dao dao.AccountDAO, cache cache.AccountCache) AccountRepository {
	return &accountRepository{dao: dao, cache: cache}
}

func (a *accountRepository) AddCredit(ctx context.Context, cr domain.Credit) error {
	activities := make([]dao.AccountActivity, 0, len(cr.Items))
	now := time.Now().UnixMilli()
	for _, itm := range cr.Items {
		activities = append(activities, dao.AccountActivity{
			Uid:         itm.Uid,
			Biz:         cr.Biz,
			BizId:       cr.BizId,
			Account:     itm.Account,
			AccountType: itm.AccountType.AsUint8(),
			Amount:      itm.Amt,
			Currency:    itm.Currency,
			Ctime:       now,
			Utime:       now,
		})
	}
	return a.dao.AddActivities(ctx, activities...)
}

func (a *accountRepository) CheckUnique(ctx context.Context, cr domain.Credit) error {
	return a.cache.GetUnique(ctx, cr)
}

func (a *accountRepository) SetUnique(ctx context.Context, cr domain.Credit) error {
	return a.cache.SetUnique(ctx, cr)
}
