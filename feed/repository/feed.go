package repository

import (
	"context"
	"encoding/json"
	"github.com/ecodeclub/ekit/slice"
	"time"
	"webook/feed/domain"
	"webook/feed/repository/dao"
)

type FeedRepository interface {
	// CreatePushEvent 批量推事件
	CreatePushEvent(ctx context.Context, fes []domain.FeedEvent) error
	// CreatePullEvent 创建拉事件
	CreatePullEvent(ctx context.Context, feedEvent domain.FeedEvent) error
	// FindPullEvents 获取拉事件，也就是关注的人发件箱里面的事件
	FindPullEvents(ctx context.Context, uids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
	// FindPushEvents 获取推事件，也就是自己收件箱里面的事件
	FindPushEvents(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
	FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
	FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
}

type feedEventRepo struct {
	pullDao dao.FeedPullEventDAO
	pushDao dao.FeedPushEventDAO
}

func NewFeedEventRepo(pullDao dao.FeedPullEventDAO, pushDao dao.FeedPushEventDAO) FeedRepository {
	return &feedEventRepo{pullDao: pullDao, pushDao: pushDao}
}

func (f *feedEventRepo) FindPullEvents(ctx context.Context, uids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pullDao.FindPullEventList(ctx, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(events, func(idx int, src dao.FeedPullEvent) domain.FeedEvent {
		return convertToPullEventDomain(src)
	}), nil
}

func (f *feedEventRepo) FindPushEvents(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pushDao.FindPushEventList(ctx, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.FeedEvent, 0, len(events))
	for _, evt := range events {
		res = append(res, convertToPushEventDomain(evt))
	}
	return res, nil
}

func (f *feedEventRepo) FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pushDao.FindPushEventListWithTyp(ctx, typ, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.FeedEvent, 0, len(events))
	for _, evt := range events {
		res = append(res, convertToPushEventDomain(evt))
	}
	return res, nil
}

func (f *feedEventRepo) FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pullDao.FindPullEventListWithTyp(ctx, typ, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(events, func(idx int, src dao.FeedPullEvent) domain.FeedEvent {
		return convertToPullEventDomain(src)
	}), nil
}

func (f *feedEventRepo) CreatePushEvent(ctx context.Context, fes []domain.FeedEvent) error {
	pushEvent := make([]dao.FeedPushEvent, 0, len(fes))
	for _, val := range fes {
		pushEvent = append(pushEvent, convertToPushEventDao(val))
	}
	return f.pushDao.CreatePushEvent(ctx, pushEvent)
}

func (f *feedEventRepo) CreatePullEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	return f.pullDao.CreatePullEvent(ctx, convertToPullEventDao(feedEvent))
}

func convertToPushEventDao(fes domain.FeedEvent) dao.FeedPushEvent {
	val, _ := json.Marshal(fes.Ext)
	return dao.FeedPushEvent{
		Uid:     fes.Uid,
		Type:    fes.Type,
		Content: string(val),
		Ctime:   fes.Ctime.UnixMilli(),
	}
}

func convertToPullEventDao(event domain.FeedEvent) dao.FeedPullEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPullEvent{
		Uid:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.UnixMilli(),
	}
}

func convertToPullEventDomain(event dao.FeedPullEvent) domain.FeedEvent {
	var ext domain.ExtendFields
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		Id:    event.Id,
		Uid:   event.Uid,
		Type:  event.Type,
		Ctime: time.UnixMilli(event.Ctime),
		Ext:   ext,
	}
}

func convertToPushEventDomain(evt dao.FeedPushEvent) domain.FeedEvent {
	var ext domain.ExtendFields
	_ = json.Unmarshal([]byte(evt.Content), &ext)
	return domain.FeedEvent{
		Id:    evt.Id,
		Uid:   evt.Uid,
		Type:  evt.Type,
		Ctime: time.UnixMilli(evt.Ctime),
		Ext:   ext,
	}
}
