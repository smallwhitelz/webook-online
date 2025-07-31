package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
	"webook/payment/domain"
	"webook/payment/repository/dao"
)

var _ LocalMsgRepository = (*LocalMsgGORMRepository)(nil)

type LocalMsgRepository interface {
	AddMsg(ctx context.Context, cnt string) (int64, error)
	MarkSuccess(ctx context.Context, msgId int64) error
	MarkFailed(ctx context.Context, msgId int64) error
	FindInitMsg(ctx context.Context, offset int, limit int) ([]domain.LocalMsg, error)
}

type LocalMsgGORMRepository struct {
	dao dao.LocalMsgDAO
}

func (l *LocalMsgGORMRepository) FindInitMsg(ctx context.Context, offset int, limit int) ([]domain.LocalMsg, error) {
	msgs, err := l.dao.FindInitMsg(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(msgs, func(idx int, src dao.LocalMsg) domain.LocalMsg {
		return domain.LocalMsg{
			Id:      src.Id,
			Content: src.Content,
			Ctime:   time.UnixMilli(src.Ctime),
			Utime:   time.UnixMilli(src.Utime),
		}
	}), nil
}

func (l *LocalMsgGORMRepository) AddMsg(ctx context.Context, cnt string) (int64, error) {
	return l.dao.AddMsg(ctx, dao.LocalMsg{
		Content: cnt,
		Status:  dao.MsgStatusInit,
	})
}

func (l *LocalMsgGORMRepository) MarkSuccess(ctx context.Context, msgId int64) error {
	return l.dao.UpdateStatus(ctx, msgId, dao.MsgStatusSuccess)
}

func (l *LocalMsgGORMRepository) MarkFailed(ctx context.Context, msgId int64) error {
	return l.dao.UpdateStatus(ctx, msgId, dao.MsgStatusFailed)
}

func NewLocalMsgGORMRepository(db *gorm.DB) LocalMsgRepository {
	return &LocalMsgGORMRepository{dao: dao.NewLocalMsgGORMDAO(db)}
}
