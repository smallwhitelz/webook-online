package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"
	"webook/interactive/domain"
	"webook/interactive/events"
	"webook/interactive/repository"
	"webook/pkg/logger"
)

//go:generate mockgen -source=./interactive.go -package=svcmocks -destination=./mocks/interactive.mock.go InteractiveService
type InteractiveService interface {
	// IncrReadCnt 增加阅读数
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	// Collect 收藏
	Collect(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
	// Get 获取文章点赞数收藏数阅读数
	Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo     repository.InteractiveRepository
	producer events.Producer
	l        logger.LoggerV1
}

func NewInteractiveService(repo repository.InteractiveRepository, producer events.Producer, l logger.LoggerV1) InteractiveService {
	return &interactiveService{repo: repo, producer: producer, l: l}
}

func (i *interactiveService) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	intrs, err := i.repo.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive, len(intrs))
	for _, intr := range intrs {
		res[intr.BizId] = intr
	}
	return res, nil
}

func (i *interactiveService) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	intr, err := i.repo.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		// 判断当前用户有没有对该文章点赞
		intr.Liked, er = i.repo.Liked(ctx, biz, bizId, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		// 判断当前用户有没有对该文章收藏
		intr.Collected, er = i.repo.Collected(ctx, biz, bizId, uid)
		return er
	})
	return intr, eg.Wait()
}

func (i *interactiveService) Collect(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error {
	err := i.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
	if err != nil {
		return err
	}
	go func() {
		err := i.sync(events.CollectEventType, biz, bizId, uid)
		if err != nil {
			i.l.Error("同步收藏数失败",
				logger.String("biz", biz),
				logger.Int64("bizId", bizId),
				logger.Error(err))
		}
	}()
	return nil
}

func (i *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	err := i.repo.IncrLike(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	go func() {
		err := i.sync(events.LikeEventType, biz, id, uid)
		if err != nil {
			i.l.Error("同步点赞数失败",
				logger.String("biz", biz),
				logger.Int64("bizId", id),
				logger.Error(err))
		}
	}()
	return nil
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := i.repo.DecrLike(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	go func() {
		err := i.sync(events.CancelLikeEventType, biz, id, uid)
		if err != nil {
			i.l.Error("同步取消点赞数失败",
				logger.String("biz", biz),
				logger.Int64("bizId", id),
				logger.Error(err))
		}
	}()
	return nil
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}

func (i *interactiveService) sync(typ events.InteractiveEventType, biz string, bizId int64, uid int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := i.producer.ProduceInteractiveEvent(ctx, events.InteractiveEvent{
		Type:  typ,
		Biz:   biz,
		BizId: bizId,
		Uid:   uid,
	})
	return err
}
