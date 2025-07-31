package service

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
	"time"
	followv1 "webook/api/proto/gen/follow/v1"
	"webook/feed/domain"
	"webook/feed/repository"
)

const (
	ArticleEventName = "article_event"
	threshold        = 4
)

type ArticleEventHandler struct {
	repo         repository.FeedRepository
	followClient followv1.FollowServiceClient
}

func NewArticleEventHandler(repo repository.FeedRepository, followClient followv1.FollowServiceClient) Handler {
	return &ArticleEventHandler{repo: repo, followClient: followClient}
}

func (a *ArticleEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	// article这边需要聚合
	// 可能在 push event 也可能在 pull event
	var eg errgroup.Group
	// 因为并发操作
	var lock sync.Mutex
	events := make([]domain.FeedEvent, 0, limit*2)
	eg.Go(func() error {
		// 查询发件箱
		// 查询我关注了那些人
		resp, err := a.followClient.GetFollowee(ctx, &followv1.GetFolloweeRequest{Follower: uid, Limit: 10000})
		if err != nil {
			return err
		}
		followeeIds := slice.Map(resp.FollowRelations, func(idx int, src *followv1.FollowRelation) int64 {
			return src.Followee
		})
		evts, err := a.repo.FindPullEventsWithTyp(ctx, ArticleEventName, followeeIds, timestamp, limit)
		if err != nil {
			return err
		}
		lock.Lock()
		events = append(events, evts...)
		lock.Unlock()
		return nil
	})
	eg.Go(func() error {
		evts, err := a.repo.FindPushEventsWithTyp(ctx, ArticleEventName, uid, timestamp, limit)
		if err != nil {
			return err
		}
		lock.Lock()
		events = append(events, evts...)
		lock.Unlock()
		return nil
	})
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	// 你已经查询到了所有数据，现在要排序
	// 这里是时间越晚越要排在前面
	// 时间越晚相当于越新
	sort.Slice(events, func(i, j int) bool {
		return events[i].Ctime.UnixMilli() > events[j].Ctime.UnixMilli()
	})
	return events[:min(int(limit), len(events))], nil
}

func (a *ArticleEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	uid, err := ext.Get("followee").AsInt64()
	if err != nil {
		return err
	}
	// 找到这个人的粉丝数量，判断是推模型还是拉模型
	resp, err := a.followClient.GetFollowStatic(ctx, &followv1.GetFollowStaticRequest{
		Uid: uid,
	})
	if err != nil {
		return err
	}
	// 大于一个阈值对应百万up主
	// 我发送一个动态，如果用写扩散，我一个人要给百万粉丝发送，这是一个很大的工作量
	// 但是我放到我的发件箱，粉丝每个人来从我这里获取信息，每个人的量是很小的，所以用读扩散
	if resp.FollowStatic.Followers > threshold {
		// 拉模型
		return a.repo.CreatePullEvent(ctx, domain.FeedEvent{
			Uid:   uid,
			Type:  ArticleEventName,
			Ctime: time.Now(),
			Ext:   ext,
		})
	} else {
		// 推模型,也就是写扩散
		// 先查询出来粉丝
		// 这里limit可以自由取值，我认为他最多只有10000个粉丝
		fresp, err := a.followClient.GetFollower(ctx, &followv1.GetFollowerRequest{
			Followee: uid,
			Limit:    10000,
			Offset:   0,
		})
		if err != nil {
			return err
		}
		events := slice.Map(fresp.FollowRelation, func(idx int, src *followv1.FollowRelation) domain.FeedEvent {
			return domain.FeedEvent{
				Uid:   src.Follower,
				Type:  ArticleEventName,
				Ctime: time.Now(),
				Ext:   ext,
			}
		})
		return a.repo.CreatePushEvent(ctx, events)
	}
}

// CreateFeedEventV1
// 如果一个用户是活跃用户，那么就保证对他进行写扩散；否则就是走原来的逻辑。
func (a *ArticleEventHandler) CreateFeedEventV1(ctx context.Context, ext domain.ExtendFields) error {
	uid, err := ext.Get("followee").AsInt64()
	if err != nil {
		return err
	}
	// 找到这个人的粉丝数量，判断是推模型还是拉模型
	resp, err := a.followClient.GetFollowStatic(ctx, &followv1.GetFollowStaticRequest{
		Uid: uid,
	})
	if err != nil {
		return err
	}
	// 大于一个阈值对应百万up主
	// 我发送一个动态，如果用写扩散，我一个人要给百万粉丝发送，这是一个很大的工作量
	// 但是我放到我的发件箱，粉丝每个人来从我这里获取信息，每个人的量是很小的，所以用读扩散
	if resp.FollowStatic.Followers > threshold {
		// 先查询出来粉丝
		// 这里limit可以自由取值，我认为他最多只有10000个粉丝
		fresp, err := a.followClient.GetFollower(ctx, &followv1.GetFollowerRequest{
			Followee: uid,
			Limit:    10000,
		})
		if err != nil {
			return err
		}
		// 粉丝里面有活跃粉丝，还是要执行写扩散，但是是针对这个用户执行写扩散
		events := slice.FilterMap(fresp.FollowRelation, func(idx int, src *followv1.FollowRelation) (domain.FeedEvent, bool) {
			// 不是活跃粉丝，什么都不做
			if !a.isActiveUser(src.Follower) {
				return domain.FeedEvent{}, false
			}
			// 活跃粉丝
			return domain.FeedEvent{Uid: src.Follower, Ctime: time.Now(), Type: ArticleEventName, Ext: ext}, true
		})
		// 活跃的话就放到粉丝的收件箱里去
		err = a.repo.CreatePushEvent(ctx, events)
		if err != nil {
			return err
		}
		// 拉模型
		return a.repo.CreatePullEvent(ctx, domain.FeedEvent{
			Uid:   uid,
			Type:  ArticleEventName,
			Ctime: time.Now(),
			Ext:   ext,
		})
	} else {
		// 先查询粉丝，而后看粉丝里面是不是有活跃用户
		fresp, err := a.followClient.GetFollower(ctx, &followv1.GetFollowerRequest{Followee: uid})
		if err != nil {
			return err
		}
		events := slice.Map(fresp.FollowRelation, func(idx int, src *followv1.FollowRelation) domain.FeedEvent {
			return domain.FeedEvent{Uid: src.Follower, Ctime: time.Now(), Type: ArticleEventName, Ext: ext}
		})
		return a.repo.CreatePushEvent(ctx, events)
	}
}

func (a *ArticleEventHandler) isActiveUser(follower int64) bool {
	return false
}
