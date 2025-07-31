package dao

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type TagDAO interface {
	CreateTag(ctx context.Context, tag Tag) (int64, error)
	GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error)
	GetTags(ctx context.Context, offset int, limit int) ([]Tag, error)
	CreateTagBiz(ctx context.Context, tagBizs []TagBiz) error
	GetTagsByBiz(ctx context.Context, uid int64, biz string, bizId int64) ([]Tag, error)
	GetTagsById(ctx context.Context, tagIds []int64) ([]Tag, error)
}

type Tag struct {
	Id   int64  `gorm:"primaryKey,autoIncrement"`
	Name string `gorm:"type=varchar(4096)"`
	// 有一个典型的场景，查处一个人有什么标签
	Uid   int64 `gorm:"index"`
	Ctime int64
	Utime int64
}

type TagBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	BizId int64  `gorm:"index:biz_type_id"`
	Biz   string `gorm:"index:biz_type_id"`
	// 冗余字段目标只有一个就是加速，加快查询和删除
	// 但是冗余字段如果被修改了，容易引起数据不一致性
	Uid   int64 `gorm:"index"`
	Tid   int64
	Tag   *Tag `gorm:"ForeignKey:Tid;AssociationForeignKey:Id;constraint:OnDelete:CASCADE"`
	Ctime int64
	Utime int64
}

type TagGORMDAO struct {
	db *gorm.DB
}

func NewTagGORMDAO(db *gorm.DB) TagDAO {
	return &TagGORMDAO{db: db}
}

func (t *TagGORMDAO) GetTagsById(ctx context.Context, tagIds []int64) ([]Tag, error) {
	var res []Tag
	err := t.db.WithContext(ctx).Where("id IN ?", tagIds).Find(&res).Error
	return res, err
}

func (t *TagGORMDAO) GetTagsByBiz(ctx context.Context, uid int64, biz string, bizId int64) ([]Tag, error) {
	// 这边使用join查询，如果你不想用join
	// 就在repo层里分为两次去查询
	// 第一次查询
	//var bizTags []TagBiz
	//err := t.db.WithContext(ctx).
	//	Where("uid = ? AND biz = ? AND biz_id = ?", uid, biz, bizId).Find(&bizTags).Error
	//if err != nil {
	//	return nil, err
	//}
	//// 第二次查询
	//ids := slice.Map(bizTags, func(idx int, src TagBiz) int64 {
	//	return src.Tid
	//})
	//var res []Tag
	//err = t.db.WithContext(ctx).Where("id IN ?", ids).Find(&res).Error
	//return res, err

	// 如果允许用join
	var bizTags []TagBiz
	err := t.db.WithContext(ctx).Model(&TagBiz{}).
		InnerJoins("Tag", t.db.Model(&Tag{})).
		Where("Tag.uid = ? AND biz = ? AND biz_id = ?", uid, biz, bizId).Find(&bizTags).Error
	if err != nil {
		return nil, err
	}
	return slice.Map(bizTags, func(idx int, src TagBiz) Tag {
		return *src.Tag
	}), nil
}

func (t *TagGORMDAO) CreateTagBiz(ctx context.Context, tagBizs []TagBiz) error {
	if len(tagBizs) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()
	for _, tb := range tagBizs {
		tb.Ctime = now
		tb.Utime = now
	}
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		first := tagBizs[0]
		err := tx.Model(&TagBiz{}).
			Delete("uid = ? AND biz = ? AND biz_id = ?", first.Uid, first.Biz, first.BizId).Error
		if err != nil {
			return err
		}
		return tx.Create(&tagBizs).Error
	})
}

func (t *TagGORMDAO) GetTags(ctx context.Context, offset int, limit int) ([]Tag, error) {
	var res []Tag
	err := t.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (t *TagGORMDAO) GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error) {
	var res []Tag
	err := t.db.WithContext(ctx).Where("uid = ?", uid).Find(&res).Error
	return res, err
}

func (t *TagGORMDAO) CreateTag(ctx context.Context, tag Tag) (int64, error) {
	now := time.Now().UnixMilli()
	tag.Ctime = now
	tag.Utime = now
	err := t.db.WithContext(ctx).Create(&tag).Error
	return tag.Id, err
}
