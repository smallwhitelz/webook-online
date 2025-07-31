package dao

import (
	"context"
	"gorm.io/gorm"
)

// FeedPullEvent 拉模型对应发件箱对应读扩散
// 你是一个百万up主，你发一个动态，你的粉丝都会从你这里去拉取你发的信息
type FeedPullEvent struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 发件人
	Uid  int64 `gorm:"index"`
	Type string
	// 这边放的就是关键的扩展字段，不同的事件类型，有不同的解析方式
	Content string
	Ctime   int64
}

type FeedPullEventDAO interface {
	CreatePullEvent(ctx context.Context, feedPullEvent FeedPullEvent) error
	FindPullEventList(ctx context.Context, uids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error)
	FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error)
}
type FeedPullEventGORMDAO struct {
	db *gorm.DB
}

func NewFeedPullEventGORMDAO(db *gorm.DB) FeedPullEventDAO {
	return &FeedPullEventGORMDAO{db: db}
}

func (f *FeedPullEventGORMDAO) FindPullEventList(ctx context.Context, uids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error) {
	var res []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("uid IN ? AND ctime < ?", uids, timestamp).
		Order("ctime desc").Limit(int(limit)).Find(&res).Error
	return res, err
}

func (f *FeedPullEventGORMDAO) FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error) {
	var res []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("uid IN ?", uids).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").Limit(int(limit)).Find(&res).Error
	return res, err
}

func (f *FeedPullEventGORMDAO) CreatePullEvent(ctx context.Context, feedPullEvent FeedPullEvent) error {
	return f.db.WithContext(ctx).Create(&feedPullEvent).Error
}
