package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/search/v1"
	"webook/search/domain"
	"webook/search/service"
)

type SyncServiceServer struct {
	searchv1.UnimplementedSyncServiceServer
	svc service.SyncService
}

func NewSyncServiceServer(svc service.SyncService) *SyncServiceServer {
	return &SyncServiceServer{svc: svc}
}

func (s *SyncServiceServer) Register(server *grpc.Server) {
	searchv1.RegisterSyncServiceServer(server, s)
}

func (s *SyncServiceServer) InputAny(ctx context.Context, request *searchv1.InputAnyRequest) (*searchv1.InputAnyResponse, error) {
	err := s.svc.InputAny(ctx, request.GetIndexName(), request.GetDocId(), request.GetData())
	return &searchv1.InputAnyResponse{}, err
}

func (s *SyncServiceServer) InputUser(ctx context.Context, request *searchv1.InputUserRequest) (*searchv1.InputUserResponse, error) {
	err := s.svc.InputUser(ctx, s.toDomainUser(request.GetUser()))
	return &searchv1.InputUserResponse{}, err
}

func (s *SyncServiceServer) InputArticle(ctx context.Context, request *searchv1.InputArticleRequest) (*searchv1.InputArticleResponse, error) {
	err := s.svc.InputArticle(ctx, s.toDomainArt(request.GetArticle()))
	return &searchv1.InputArticleResponse{}, err
}

func (s *SyncServiceServer) toDomainUser(user *searchv1.User) domain.User {
	return domain.User{
		Id:       user.GetId(),
		Nickname: user.GetNickname(),
		Email:    user.GetEmail(),
		Phone:    user.GetPhone(),
	}
}

func (s *SyncServiceServer) toDomainArt(article *searchv1.Article) domain.Article {
	return domain.Article{
		Id:      article.GetId(),
		Title:   article.GetTitle(),
		Content: article.GetContent(),
		Status:  article.GetStatus(),
	}
}
