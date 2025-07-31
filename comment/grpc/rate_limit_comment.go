package grpc

import (
	"context"
	"errors"
	commentv1 "webook/api/proto/gen/comment/v1"
)

type RateLimitComment struct {
	CommentServiceServer
}

// GetCommentList 保护系统的玩法
// 这种就是以热门资源和非热门资源的玩法进行限流
func (r *RateLimitComment) GetCommentList(ctx context.Context, request *commentv1.GetCommentListRequest) (*commentv1.GetCommentListResponse, error) {
	// 一般是通过热榜功能，提前计算放到redis里，问一下redis就知道是不是热门资源了
	isHotBiz := r.isHotBiz(request.GetBiz(), request.GetBizId())
	if isHotBiz {
		// 限流阈值 400/s
	} else {
		// 限流阈值 100/s
	}
	return r.CommentServiceServer.GetCommentList(ctx, request)
}

// GetCommentListV1
// 这种玩法是限流并且降级了
// 这个时候我们就只让他查询热门资源
func (r *RateLimitComment) GetCommentListV1(ctx context.Context, request *commentv1.GetCommentListRequest) (*commentv1.GetCommentListResponse, error) {
	// 一般是通过热榜功能，提前计算放到redis里，问一下redis就知道是不是热门资源了
	isHotBiz := r.isHotBiz(request.GetBiz(), request.GetBizId())
	if !isHotBiz && ctx.Value("downgrade") == "true" {
		// 非热门资源触发降级
		return &commentv1.GetCommentListResponse{}, errors.New("非热门资源降级")
	} else {
		// 限流阈值 100/s
	}
	return r.CommentServiceServer.GetCommentList(ctx, request)
}

// CreateComment 写评论这里还可以加入异步写
// 比如你是一个百万up主，一发文章就有很多人给你评论
//  1. 这个时候就会达到峰值，达到峰值后我们可以异步借助kafka去消费写到数据库，进行流量削峰
//  2. 还有一种就是现在基本都是先审后发的情况，所以你发完评论，可以让审核那边进行消费，如果没有通过
//     那么这条评论只有你自己可以看到，并不会公开
func (r *RateLimitComment) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	if ctx.Value("limited") == "true" || ctx.Value("downgrade") == "true" {
		// 转kafka
		return &commentv1.CreateCommentResponse{}, nil
	}
	err := r.svc.CreateComment(ctx, r.convertToDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}

func (r *RateLimitComment) isHotBiz(biz string, bizId int64) bool {
	return true
}
