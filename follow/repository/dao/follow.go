package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type FollowGROMDAO struct {
	db *gorm.DB
}

func NewFollowGROMDAO(db *gorm.DB) FollowDAO {
	return &FollowGROMDAO{db: db}
}

func (f *FollowGROMDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := f.db.WithContext(ctx).Select("count(follower)").
		// 如果没有额外的索引，绝对是全表扫描
		// 可以考虑在followee上建立一个额外的索引
		Where("followee = ? AND status = ?", uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (f *FollowGROMDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := f.db.WithContext(ctx).Select("count(followee)").
		// <follower,followee> 会命中索引
		Where("follower = ? AND status = ?", uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (f *FollowGROMDAO) FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error) {
	var res FollowRelation
	err := f.db.WithContext(ctx).
		// 这里感觉有待考虑
		// 我们这个接口只是查询关注的人进去的状态信息吗？
		// 未关注的人的展示信息如何拿呢？
		Where("follower = ? AND followee = ? AND status = ?", follower, followee, FollowRelationStatusActive).
		First(&res).Error
	return res, err
}

func (f *FollowGROMDAO) GetFolloweeList(ctx context.Context, follower int64, offset int64, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := f.db.WithContext(ctx).
		Where("follower = ? AND status = ?", follower, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).Find(&res).Error
	return res, err
}

func (f *FollowGROMDAO) GetFollowerList(ctx context.Context, followee int64, offset int64, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := f.db.WithContext(ctx).
		Where("followee = ? AND status = ?", followee, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).Find(&res).Error
	return res, err
}

func (f *FollowGROMDAO) UpdateStatus(ctx context.Context, fr FollowRelation, status uint8) error {
	// 如果当前status就是inactive呢？
	// 不用去管，也不用去检测，正常人不会取消关注还继续调取消关注的接口
	// 黑客的体验感用管吗？不用
	return f.db.WithContext(ctx).
		Where("follower = ? AND followee = ?", fr.Follower, fr.Followee).
		Updates(map[string]any{
			"status": status,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (f *FollowGROMDAO) CreateFollowRelation(ctx context.Context, fr FollowRelation) error {
	// 我也要保持insert or update语义
	// 因为这里要保证你点完关注再点一定就是取消关注
	now := time.Now().UnixMilli()
	fr.Ctime = now
	fr.Utime = now
	fr.Status = FollowRelationStatusActive
	return f.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			// 代表关注-取消-再关注
			"status": FollowRelationStatusActive,
			"utime":  now,
		}),
	}).Create(&fr).Error
}
