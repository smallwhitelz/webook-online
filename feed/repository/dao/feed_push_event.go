package dao

import (
	"context"
	"gorm.io/gorm"
)

// FeedPushEvent 推模型对应收件箱对应写扩散
// A 关注了 B，那么 B 在发表一篇新的文章的时候，就会把这条数据写入到 A 的收件箱。
type FeedPushEvent struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 收件人
	Uid  int64 `gorm:"index"`
	Type string
	// 这边放的就是关键的扩展字段，不同的事件类型，有不同的解析方式
	Content string
	Ctime   int64
	// 正常来说，这个表的数据是不会被更新的
	//Utime int64
}

type FeedPushEventDAO interface {
	CreatePushEvent(ctx context.Context, pushEvent []FeedPushEvent) error
	FindPushEventList(ctx context.Context, uid int64, timestamp int64, limit int64) ([]FeedPushEvent, error)
	FindPushEventListWithTyp(ctx context.Context, typ string, uid int64, timestamp int64, limit int64) ([]FeedPushEvent, error)
}

type FeedPushEventGORMDAO struct {
	db *gorm.DB
}

func NewFeedPushEventGORMDAO(db *gorm.DB) FeedPushEventDAO {
	return &FeedPushEventGORMDAO{db: db}
}

func (f *FeedPushEventGORMDAO) FindPushEventList(ctx context.Context, uid int64, timestamp int64, limit int64) ([]FeedPushEvent, error) {
	var res []FeedPushEvent
	err := f.db.WithContext(ctx).
		Where("uid = ? AND ctime < ?", uid, timestamp).
		Order("ctime desc").Limit(int(limit)).Find(&res).Error
	return res, err
}

func (f *FeedPushEventGORMDAO) FindPushEventListWithTyp(ctx context.Context, typ string, uid int64, timestamp int64, limit int64) ([]FeedPushEvent, error) {
	var res []FeedPushEvent
	err := f.db.WithContext(ctx).
		Where("uid = ? AND ctime < ? AND type = ?", uid, timestamp, typ).
		Order("ctime desc").Limit(int(limit)).Find(&res).Error
	return res, err
}

func (f *FeedPushEventGORMDAO) CreatePushEvent(ctx context.Context, pushEvent []FeedPushEvent) error {
	return f.db.WithContext(ctx).Create(&pushEvent).Error
}
