package service

import (
	"context"
	"webook/comment/domain"
	"webook/comment/repository"
)

type CommentService interface {
	// CreateComment 创建评论
	CreateComment(ctx context.Context, cm domain.Comment) error
	// DeleteComment 删除评论，删除本评论和其子评论
	DeleteComment(ctx context.Context, id int64) error
	// GetCommentList Comment的id为0 获取一级评论
	// 按照 ID 倒序排序
	GetCommentList(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error)
	// GetMoreReplies 查询更多评论，类似bilibili的评论点击查看更多
	// 升序
	GetMoreReplies(ctx context.Context, rid int64, maxId int64, limit int64) ([]domain.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (c *commentService) GetMoreReplies(ctx context.Context, rid int64, maxId int64, limit int64) ([]domain.Comment, error) {
	list, err := c.repo.GetMoreReplies(ctx, rid, maxId, limit)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// CreateComment 写评论这里还可以加入异步写
// 比如你是一个百万up主，一发文章就有很多人给你评论
//  1. 这个时候就会达到峰值，达到峰值后我们可以异步借助kafka去消费写到数据库，进行流量削峰
//  2. 还有一种就是现在基本都是先审后发的情况，所以你发完评论，可以让审核那边进行消费，如果没有通过
//     那么这条评论只有你自己可以看到，并不会公开
func (c *commentService) CreateComment(ctx context.Context, cm domain.Comment) error {
	return c.repo.CreateComment(ctx, cm)
}

// DeleteComment 这里我们认为父评论删除了，子评论也应该被删除，避免出现孤儿节点
func (c *commentService) DeleteComment(ctx context.Context, id int64) error {
	return c.repo.DeleteComment(ctx, domain.Comment{
		Id: id,
	})
}

func (c *commentService) GetCommentList(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error) {
	list, err := c.repo.FindCommentByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	return list, nil
}
