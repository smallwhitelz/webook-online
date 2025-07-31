package service

import (
	"context"
	"webook/feed/domain"
)

type FeedService interface {
	CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error
	GetFeedEventList(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error)
}

// Handler 主要处理和具体业务有关的内容
// 而且这里也遵循了开闭原则
// 新业务要接入进来，实现该接口即可
// 负责判断要用推模型还是拉模型
// 可以检测业务方传入的数据是否缺少扩展字段
// 如果业务方没有传入一些在展示用的字段，那么可以在这里查询，例如根据传过来的uid找到对应user的nickname
type Handler interface {
	CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error
	FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error)
}
