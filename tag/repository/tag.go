package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"time"
	"webook/pkg/logger"
	"webook/tag/domain"
	"webook/tag/repository/cache"
	"webook/tag/repository/dao"
)

type TagRepository interface {
	CreateTag(ctx context.Context, tag domain.Tag) (int64, error)
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tagIds []int64) error
	GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error)
	GetTagsById(ctx context.Context, tagIds []int64) ([]domain.Tag, error)
}

type tagRepository struct {
	dao   dao.TagDAO
	cache cache.TagCache
	l     logger.LoggerV1
}

func NewTagRepository(dao dao.TagDAO, cache cache.TagCache, l logger.LoggerV1) TagRepository {
	return &tagRepository{dao: dao, cache: cache, l: l}
}

func (t *tagRepository) GetTagsById(ctx context.Context, tagIds []int64) ([]domain.Tag, error) {
	tags, err := t.dao.GetTagsById(ctx, tagIds)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return t.toDomain(src)
	}), nil
}

func (t *tagRepository) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	tags, err := t.dao.GetTagsByBiz(ctx, uid, biz, bizId)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return t.toDomain(src)
	}), nil
}

func (t *tagRepository) BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tagIds []int64) error {
	//return t.dao.CreateTagBiz(ctx, slice.Map(tagIds, func(idx int, src int64) dao.TagBiz {
	//	return dao.TagBiz{
	//		Tid:   src,
	//		BizId: bizId,
	//		Biz:   biz,
	//		Uid:   uid,
	//	}
	//}))
	tagBizList := make([]dao.TagBiz, 0, len(tagIds))
	for _, id := range tagIds {
		tagBizList = append(tagBizList, dao.TagBiz{
			BizId: bizId,
			Biz:   biz,
			Uid:   uid,
			Tid:   id,
		})
	}
	return t.dao.CreateTagBiz(ctx, tagBizList)
}

// PreloadUsersTags 在toB的场景下，可以提前预加载缓存
// 用在ioc里，也就是这次在初始化要单独初始化repository了
func (t *tagRepository) PreloadUsersTags(ctx context.Context) error {
	// 我们要存的是uid => 我的所有标签
	// 这边分批预加载出来
	// 数据从里面取出来，调用append
	offset := 0
	const limit = 100
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		tags, err := t.dao.GetTags(dbCtx, offset, limit)
		cancel()
		if err != nil {
			// 也可以continue
			return err
		}
		for _, tag := range tags {
			rcCtx, cancel := context.WithTimeout(ctx, time.Second)
			err := t.cache.Append(rcCtx, tag.Uid, t.toDomain(tag))
			cancel()
			if err != nil {
				continue
			}
		}
		if len(tags) < limit {
			return nil
		}
		offset = offset + limit
	}
}

func (t *tagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	res, err := t.cache.GetTags(ctx, uid)
	if err == nil {
		return res, nil
	}
	tags, err := t.dao.GetTagsByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	res = slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return t.toDomain(src)
	})
	err = t.cache.Append(ctx, uid, res...)
	if err != nil {
		// 记录日志
		t.l.Error("获取tags放入缓存失败",
			logger.Error(err),
			logger.Int64("uid", uid))
	}
	return res, nil
}

func (t *tagRepository) CreateTag(ctx context.Context, tag domain.Tag) (int64, error) {
	id, err := t.dao.CreateTag(ctx, t.toEntity(tag))
	if err != nil {
		return 0, err
	}
	err = t.cache.Append(ctx, tag.Uid, tag)
	if err != nil {
		// 记录日志即可
		t.l.Error("创建tag放入缓存中失败", logger.Error(err),
			logger.Int64("uid", tag.Uid), logger.String("name", tag.Name))
	}
	return id, nil
}

func (t *tagRepository) toDomain(src dao.Tag) domain.Tag {
	return domain.Tag{
		Id:    src.Id,
		Name:  src.Name,
		Uid:   src.Uid,
		Ctime: time.UnixMilli(src.Ctime),
		Utime: time.UnixMilli(src.Utime),
	}
}

func (t *tagRepository) toEntity(tag domain.Tag) dao.Tag {
	return dao.Tag{
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}
