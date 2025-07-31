package service

import (
	"context"
	"time"
	"webook/feed/domain"
	"webook/feed/repository"
)

const (
	LikeEventName = "like_event"
)

// LikeEventHandler 点赞这件事
// 只有点赞人和被点赞人两个人知道，所以这里用推模型
type LikeEventHandler struct {
	repo repository.FeedRepository
}

func NewLikeEventHandler(repo repository.FeedRepository) Handler {
	return &LikeEventHandler{repo: repo}
}

func (l *LikeEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	return l.repo.FindPushEventsWithTyp(ctx, LikeEventName, uid, timestamp, limit)
}

// CreateFeedEvent 中的ext至少需要有三个id
// liked int64：被点赞的人
// liker int64：点赞人
// bizId int64：被点赞的东西
// biz string：业务
func (l *LikeEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	// 你可以在这里校验字段
	// 我现在需要被点赞的人，因为要放到被点赞的人的收件箱里去
	uid, err := ext.Get("liked").AsInt64()
	if err != nil {
		return err
	}
	// 比如说你准备冗余，但是业务方又不愿意提供冗余字段
	// 你可以去查，比如说在这里调用 biz + biz id 拿到资源数据
	// 调用 user 拿到用户昵称
	return l.repo.CreatePushEvent(ctx, []domain.FeedEvent{
		{Uid: uid, Ext: ext, Type: LikeEventName, Ctime: time.Now()},
	})
}
