package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"webook/api/proto/gen/comment/v1"
	"webook/comment/domain"
	"webook/comment/service"
)

type CommentServiceServer struct {
	commentv1.UnimplementedCommentServiceServer
	svc service.CommentService
}

func NewCommentServiceServer(svc service.CommentService) *CommentServiceServer {
	return &CommentServiceServer{svc: svc}
}

func (c *CommentServiceServer) Register(s *grpc.Server) {
	commentv1.RegisterCommentServiceServer(s, c)
}

func (c *CommentServiceServer) GetMoreReplies(ctx context.Context, request *commentv1.GetMoreRepliesRequest) (*commentv1.GetMoreRepliesResponse, error) {
	dmComments, err := c.svc.GetMoreReplies(ctx, request.GetRid(), request.GetMaxId(), request.GetLimit())
	if err != nil {
		return nil, err
	}
	return &commentv1.GetMoreRepliesResponse{
		Replies: c.toDTO(dmComments),
	}, nil
}

func (c *CommentServiceServer) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	err := c.svc.CreateComment(ctx, c.convertToDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}

func (c *CommentServiceServer) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	err := c.svc.DeleteComment(ctx, request.GetId())
	return &commentv1.DeleteCommentResponse{}, err
}

func (c *CommentServiceServer) GetCommentList(ctx context.Context, request *commentv1.GetCommentListRequest) (*commentv1.GetCommentListResponse, error) {
	minId := request.GetMinId()
	// 第一次查询
	if minId <= 0 {
		minId = math.MaxInt64
	}
	dmComments, err := c.svc.GetCommentList(ctx, request.GetBiz(), request.GetBizId(), minId, request.GetLimit())
	if err != nil {
		return nil, err
	}
	return &commentv1.GetCommentListResponse{
		Comments: c.toDTO(dmComments),
	}, nil
}

func (c *CommentServiceServer) convertToDomain(comment *commentv1.Comment) domain.Comment {
	dc := domain.Comment{
		Id: comment.Id,
		Commentator: domain.User{
			Id: comment.Uid,
		},
		Biz:     comment.Biz,
		BizId:   comment.BizId,
		Content: comment.Content,
		Ctime:   comment.GetCtime().AsTime(),
		Utime:   comment.GetUtime().AsTime(),
	}
	if comment.GetRootComment() != nil {
		dc.RootComment = &domain.Comment{
			Id: comment.GetRootComment().GetId(),
		}
	}
	if comment.GetParentComment() != nil {
		dc.ParentComment = &domain.Comment{
			Id: comment.GetParentComment().GetId(),
		}
	}
	return dc
}

func (c *CommentServiceServer) toDTO(dmComments []domain.Comment) []*commentv1.Comment {
	rpcComments := make([]*commentv1.Comment, 0, len(dmComments))
	for _, cmt := range dmComments {
		rpcComment := &commentv1.Comment{
			Id:      cmt.Id,
			Uid:     cmt.Commentator.Id,
			Biz:     cmt.Biz,
			BizId:   cmt.BizId,
			Content: cmt.Content,
			Ctime:   timestamppb.New(cmt.Ctime),
			Utime:   timestamppb.New(cmt.Utime),
		}
		if cmt.RootComment != nil {
			rpcComment.RootComment = &commentv1.Comment{
				Id: cmt.RootComment.Id,
			}
		}
		if cmt.ParentComment != nil {
			rpcComment.ParentComment = &commentv1.Comment{
				Id: cmt.ParentComment.Id,
			}
		}
		rpcComments = append(rpcComments, rpcComment)
	}
	rpcCommentMap := make(map[int64]*commentv1.Comment, len(rpcComments))
	for _, rpcComment := range rpcComments {
		rpcCommentMap[rpcComment.Id] = rpcComment
	}
	for _, dmComment := range dmComments {
		rpcComment := rpcCommentMap[dmComment.Id]
		if dmComment.RootComment != nil {
			val, ok := rpcCommentMap[dmComment.RootComment.Id]
			if ok {
				rpcComment.RootComment = val
			}
		}
		if dmComment.ParentComment != nil {
			val, ok := rpcCommentMap[dmComment.ParentComment.Id]
			if ok {
				rpcComment.ParentComment = val
			}
		}
	}
	return rpcComments
}
