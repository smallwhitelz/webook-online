package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AccountGORMDAO struct {
	db *gorm.DB
}

func NewAccountGORMDAO(db *gorm.DB) AccountDAO {
	return &AccountGORMDAO{db: db}
}

// AddActivities 在金融行业，类似支付宝、wxzf这种，他们的账号服务在这个地方不一定可以使用本地事务，因为数据库可能撑不住
// 分库分表也不行，因为一笔账涉及的多个账号完全有可能在不同库上
// 有一种办法 可以使用钞能力，买oracle、服务器、买操作系统
func (a *AccountGORMDAO) AddActivities(ctx context.Context, activities ...AccountActivity) error {
	// 记录账户流水的时候，余额也应该进行更新，这两者应该在同一个事务里，同时失败同时成功
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 针对每一个activity 入账
		now := time.Now().UnixMilli()
		for _, act := range activities {
			// 一般在用户注册好的时候，我们就要创建好账号，所以兼容处理一下
			// 系统账号是默认一定存在的，一般是离线创建好的
			// 先更新余额
			err := tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]any{
					"balance": gorm.Expr(" `balance` + ? ", act.Amount),
					"utime":   now,
				}),
			}).Create(&Account{
				Uid:      act.Uid,
				Account:  act.Account,
				Type:     act.AccountType,
				Balance:  act.Amount,
				Currency: act.Currency,
				Ctime:    now,
				Utime:    now,
			}).Error
			if err != nil {
				return err
			}
		}
		return tx.Create(&activities).Error
	})
}
