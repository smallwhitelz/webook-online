package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"webook/api/proto/gen/ranking/v1"
	"webook/ranking/domain"
	"webook/ranking/service"
)

type RankingServiceServer struct {
	rankingv1.UnimplementedRankingServiceServer
	svc service.RankingService
}

func NewRankingServiceServer(svc service.RankingService) *RankingServiceServer {
	return &RankingServiceServer{svc: svc}
}

func (r *RankingServiceServer) Register(server *grpc.Server) {
	rankingv1.RegisterRankingServiceServer(server, r)
}

func (r *RankingServiceServer) TopN(ctx context.Context, request *rankingv1.TopNRequest) (*rankingv1.TopNResponse, error) {
	err := r.svc.TopN(ctx)
	return &rankingv1.TopNResponse{}, err
}

func (r *RankingServiceServer) GetTopN(ctx context.Context, request *rankingv1.GetTopNRequest) (*rankingv1.GetTopNResponse, error) {
	articles, err := r.svc.GetTopN(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*rankingv1.Article, 0, len(articles))
	for _, article := range articles {
		res = append(res, r.convertToView(article))
	}
	return &rankingv1.GetTopNResponse{
		Articles: res,
	}, nil
}

func (r *RankingServiceServer) convertToView(article domain.Article) *rankingv1.Article {
	return &rankingv1.Article{
		Id:      article.Id,
		Title:   article.Title,
		Status:  int32(article.Status),
		Content: article.Content,
		Author: &rankingv1.Author{
			Id:   article.Author.Id,
			Name: article.Author.Name,
		},
		Ctime: timestamppb.New(article.Ctime),
		Utime: timestamppb.New(article.Utime),
	}
}
