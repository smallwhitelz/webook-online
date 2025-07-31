package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type LocalMsg struct {
	Id      int64 `gorm:"primaryKey,autoIncrement"`
	Content string
	Status  uint8
	Ctime   int64
	Utime   int64
}

const (
	MsgStatusUnknown = 0
	MsgStatusInit    = 1
	MsgStatusSuccess = 2
	MsgStatusFailed  = 3
)

// 在编译时检查*LocalMsgGORMDAO是否实现了LocalMsgDAO接口的所有方法。如果没有实现，编译会直接失败。
var _ LocalMsgDAO = (*LocalMsgGORMDAO)(nil)

type LocalMsgDAO interface {
	AddMsg(ctx context.Context, msg LocalMsg) (int64, error)
	UpdateStatus(ctx context.Context, msgId int64, status uint8) error
	FindInitMsg(ctx context.Context, offset int, limit int) ([]LocalMsg, error)
}

type LocalMsgGORMDAO struct {
	db *gorm.DB
}

func (g *LocalMsgGORMDAO) AddMsg(ctx context.Context, msg LocalMsg) (int64, error) {
	now := time.Now().UnixMilli()
	msg.Ctime = now
	msg.Utime = now
	err := g.db.WithContext(ctx).Create(&msg).Error
	return msg.Id, err
}

func (g *LocalMsgGORMDAO) UpdateStatus(ctx context.Context, msgId int64, status uint8) error {
	return g.db.WithContext(ctx).Model(&LocalMsg{}).Where("id = ?", msgId).
		Updates(map[string]any{
			"status": status,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *LocalMsgGORMDAO) FindInitMsg(ctx context.Context, offset int, limit int) ([]LocalMsg, error) {
	var res []LocalMsg
	err := g.db.WithContext(ctx).Where("status = ?", MsgStatusInit).
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func NewLocalMsgGORMDAO(db *gorm.DB) LocalMsgDAO {
	return &LocalMsgGORMDAO{db: db}
}
