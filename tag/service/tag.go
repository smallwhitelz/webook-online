package service

import (
	"context"
	"encoding/json"
	"github.com/ecodeclub/ekit/slice"
	"time"
	"webook/pkg/logger"
	"webook/tag/domain"
	"webook/tag/events"
	"webook/tag/repository"
)

type TagService interface {
	// CreateTag 用户创建标签
	CreateTag(ctx context.Context, uid int64, name string) (int64, error)
	// GetTags 用户获取自己创建的标签
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	// AttachTags 覆盖资源标签
	AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tagIds []int64) error
	// GetBizTags 查找资源的标签
	GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error)
}

type tagService struct {
	repo     repository.TagRepository
	producer events.Producer
	l        logger.LoggerV1
}

func NewTagService(repo repository.TagRepository, producer events.Producer, l logger.LoggerV1) TagService {
	return &tagService{repo: repo, producer: producer, l: l}
}

func (t *tagService) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	return t.repo.GetBizTags(ctx, uid, biz, bizId)
}

func (t *tagService) AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tagIds []int64) error {
	err := t.repo.BindTagToBiz(ctx, uid, biz, bizId, tagIds)
	if err != nil {
		return err
	}
	// 异步发送
	// 当用户给某个资源加标签时，将该资源和标签同步到搜索里
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		tags, err := t.repo.GetTagsById(ctx, tagIds)
		if err != nil {
			return
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		// 这里要注意顺序
		// 一个用户他可能给一个资源打标签，先打abc，再打的bcd
		// 这里要注意消息发送的顺序，所以要用hash类partition
		evt := events.BizTags{
			Biz:   biz,
			BizId: bizId,
			Uid:   uid,
			Tags: slice.Map(tags, func(idx int, src domain.Tag) string {
				return src.Name
			}),
		}
		err = t.producer.ProduceSyncEvent(ctx, evt)
		cancel()
		if err != nil {
			// 记录日志即可
			val, err := json.Marshal(evt)
			if err != nil {
				return
			}
			t.l.Error("同步tags的生产者发送消息失败",
				logger.Error(err),
				logger.String("biz", biz),
				logger.Int64("biz_id", bizId),
				logger.Int64("uid", uid),
				logger.String("evt", string(val)))
		}
	}()
	return err
}

func (t *tagService) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	return t.repo.GetTags(ctx, uid)
}

func (t *tagService) CreateTag(ctx context.Context, uid int64, name string) (int64, error) {
	id, err := t.repo.CreateTag(ctx, domain.Tag{
		Uid:  uid,
		Name: name,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
