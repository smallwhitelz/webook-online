package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/search/v1"
	"webook/search/domain"
	"webook/search/service"
)

type SearchServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	svc service.SearchService
}

func NewSearchServiceServer(svc service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{svc: svc}
}

func (s *SearchServiceServer) Register(server *grpc.Server) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func (s *SearchServiceServer) Search(ctx context.Context, request *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	searchResult, err := s.svc.Search(ctx, request.GetUid(), request.GetExpression())
	if err != nil {
		return nil, err
	}
	userRes := make([]*searchv1.User, 0, len(searchResult.Users))
	artRes := make([]*searchv1.Article, 0, len(searchResult.Articles))
	for _, user := range searchResult.Users {
		userRes = append(userRes, s.toViewUser(user))
	}
	for _, article := range searchResult.Articles {
		artRes = append(artRes, s.toViewArt(article))
	}
	return &searchv1.SearchResponse{
		User: &searchv1.UserResult{
			Users: userRes,
		},
		Article: &searchv1.ArticleResult{
			Articles: artRes,
		},
	}, nil
}

func (s *SearchServiceServer) toViewUser(user domain.User) *searchv1.User {
	return &searchv1.User{
		Id:       user.Id,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
	}
}

func (s *SearchServiceServer) toViewArt(article domain.Article) *searchv1.Article {
	return &searchv1.Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Status:  article.Status,
	}
}
