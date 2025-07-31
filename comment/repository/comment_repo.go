package repository

import (
	"context"
	"database/sql"
	"golang.org/x/sync/errgroup"
	"time"
	"webook/comment/domain"
	"webook/comment/repository/dao"
	"webook/pkg/logger"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, cm domain.Comment) error
	DeleteComment(ctx context.Context, cm domain.Comment) error
	FindCommentByBiz(ctx context.Context, biz string, bizId int64, minID int64, limit int64) ([]domain.Comment, error)
	GetMoreReplies(ctx context.Context, rid int64, maxId int64, limit int64) ([]domain.Comment, error)
}

type commentRepository struct {
	dao dao.CommentDAO
	l   logger.LoggerV1
}

func NewCommentRepository(dao dao.CommentDAO, l logger.LoggerV1) CommentRepository {
	return &commentRepository{dao: dao, l: l}
}

func (c *commentRepository) GetMoreReplies(ctx context.Context, rid int64, maxId int64, limit int64) ([]domain.Comment, error) {
	cs, err := c.dao.FindRepliesByRid(ctx, rid, maxId, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Comment, 0, len(cs))
	for _, val := range cs {
		res = append(res, c.toDomain(val))
	}
	return res, nil
}

// FindCommentByBiz 有没有必要加缓存
// 这里其实要和热榜结合起来
// 假如up主是一个小菜，发布的东西几个月没人点击或者评论，那就没必要去缓存
// 所以这里的缓存一定要和热榜结合
func (c *commentRepository) FindCommentByBiz(ctx context.Context, biz string, bizId int64, minID int64, limit int64) ([]domain.Comment, error) {
	daoCmts, err := c.dao.FindCommentByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	// 这时候去查找子评论，找三条
	res := make([]domain.Comment, 0, len(daoCmts))
	// 只找三条
	var eg errgroup.Group
	// 如果触发了熔断
	//downgrade := ctx.Value("downgrade") == "true"
	for _, dc := range daoCmts {
		// 保证下面的goroutine执行时拿到的还是同一个dc
		// 没有这一步，有可能导致在这里还是dc0，到另一个goroutine里变更dc1
		dc := dc
		cm := c.toDomain(dc)
		res = append(res, cm)
		// 那么就可以不用找他的子评论
		//if downgrade {
		//	continue
		//}
		eg.Go(func() error {

			// 只展示三条
			cm.Children = make([]domain.Comment, 0, 3)
			subCmts, err := c.dao.FindRepliesByPid(ctx, dc.Id, 0, 3)
			if err != nil {
				// 我们认为这是一个可以容忍的错误
				c.l.Error("查询子评论失败", logger.Error(err))
				return nil
			}

			for _, sc := range subCmts {
				cm.Children = append(cm.Children, c.toDomain(sc))
			}
			return nil
		})
	}
	// eg.Wait() 等所有的都执行完
	return res, eg.Wait()
}

func (c *commentRepository) DeleteComment(ctx context.Context, cm domain.Comment) error {
	return c.dao.Delete(ctx, dao.Comment{
		Id: cm.Id,
	})
}

func (c *commentRepository) CreateComment(ctx context.Context, cm domain.Comment) error {
	return c.dao.Insert(ctx, c.toEntity(cm))
}

func (c *commentRepository) toEntity(cm domain.Comment) dao.Comment {
	dc := dao.Comment{
		Uid:     cm.Commentator.Id,
		Biz:     cm.Biz,
		BizId:   cm.BizId,
		Content: cm.Content,
		Ctime:   time.Now().UnixMilli(),
		Utime:   time.Now().UnixMilli(),
	}
	// 如果评论有根评论，父评论
	// 将其id赋值给对应字段
	// 否则将为null
	if cm.RootComment != nil {
		dc.RootID = sql.NullInt64{
			Valid: true,
			Int64: cm.RootComment.Id,
		}
	}
	if cm.ParentComment != nil {
		dc.PID = sql.NullInt64{
			Valid: true,
			Int64: cm.ParentComment.Id,
		}
	}
	return dc
}

func (c *commentRepository) toDomain(dc dao.Comment) domain.Comment {
	val := domain.Comment{
		Id: dc.Id,
		Commentator: domain.User{
			Id: dc.Uid,
		},
		Biz:     dc.Biz,
		BizId:   dc.BizId,
		Content: dc.Content,
		Ctime:   time.UnixMilli(dc.Ctime),
		Utime:   time.UnixMilli(dc.Utime),
	}
	if dc.RootID.Valid {
		val.RootComment = &domain.Comment{
			Id: dc.RootID.Int64,
		}
	}
	if dc.PID.Valid {
		val.ParentComment = &domain.Comment{
			Id: dc.PID.Int64,
		}
	}
	return val
}
