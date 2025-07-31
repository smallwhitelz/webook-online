package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"webook/reward/domain"
)

type RewardGORMDAO struct {
	db *gorm.DB
}

func NewRewardGORMDAO(db *gorm.DB) RewardDAO {
	return &RewardGORMDAO{db: db}
}

func (dao *RewardGORMDAO) GetReward(ctx context.Context, rid int64) (Reward, error) {
	var res Reward
	err := dao.db.WithContext(ctx).Where("id = ?").First(&res).Error
	return res, err
}

func (dao *RewardGORMDAO) UpdateReward(ctx context.Context, rid int64, status domain.RewardStatus) error {
	return dao.db.WithContext(ctx).Where("id = ?", rid).Updates(map[string]any{
		"status": status.AsUint8(),
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (dao *RewardGORMDAO) Insert(ctx context.Context, r Reward) (int64, error) {
	now := time.Now().UnixMilli()
	r.Ctime = now
	r.Utime = now
	err := dao.db.WithContext(ctx).Create(&r).Error
	return r.Id, err
}
