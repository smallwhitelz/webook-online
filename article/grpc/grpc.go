package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	articlev1 "webook/api/proto/gen/article/v1"
	"webook/article/domain"
	"webook/article/service"
)

type ArticleServiceServer struct {
	articlev1.UnimplementedArticleServiceServer
	svc service.ArticleService
}

func NewArticleServiceServer(svc service.ArticleService) *ArticleServiceServer {
	return &ArticleServiceServer{svc: svc}
}

func (a *ArticleServiceServer) Register(server *grpc.Server) {
	articlev1.RegisterArticleServiceServer(server, a)
}

func (a *ArticleServiceServer) Save(ctx context.Context, request *articlev1.SaveRequest) (*articlev1.SaveResponse, error) {
	artId, err := a.svc.Save(ctx, a.convertToDomain(request.GetArticle()))
	if err != nil {
		return nil, err
	}
	return &articlev1.SaveResponse{
		Id: artId,
	}, err
}

func (a *ArticleServiceServer) Publish(ctx context.Context, request *articlev1.PublishRequest) (*articlev1.PublishResponse, error) {
	artID, err := a.svc.Publish(ctx, a.convertToDomain(request.GetArticle()))
	if err != nil {
		return nil, err
	}
	return &articlev1.PublishResponse{
		Id: artID,
	}, nil
}

func (a *ArticleServiceServer) Withdraw(ctx context.Context, request *articlev1.WithdrawRequest) (*articlev1.WithdrawResponse, error) {
	err := a.svc.Withdraw(ctx, request.GetUid(), request.GetAid())
	return &articlev1.WithdrawResponse{}, err
}

func (a *ArticleServiceServer) GetByAuthor(ctx context.Context, request *articlev1.GetByAuthorRequest) (*articlev1.GetByAuthorResponse, error) {
	arts, err := a.svc.GetByAuthor(ctx, request.GetUid(), int(request.GetOffset()), int(request.GetLimit()))
	if err != nil {
		return nil, err
	}
	res := make([]*articlev1.Article, 0, len(arts))
	for _, art := range arts {
		res = append(res, a.convertToView(art))
	}
	return &articlev1.GetByAuthorResponse{
		Articles: res,
	}, nil
}

func (a *ArticleServiceServer) GetById(ctx context.Context, request *articlev1.GetByIdRequest) (*articlev1.GetByIdResponse, error) {
	art, err := a.svc.GetById(ctx, request.GetAid())
	if err != nil {
		return nil, err
	}
	return &articlev1.GetByIdResponse{
		Article: a.convertToView(art),
	}, nil
}

func (a *ArticleServiceServer) GetPubById(ctx context.Context, request *articlev1.GetPubByIdRequest) (*articlev1.GetPubByIdResponse, error) {
	art, err := a.svc.GetPubById(ctx, request.GetAid(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &articlev1.GetPubByIdResponse{
		Article: a.convertToView(art),
	}, nil
}

func (a *ArticleServiceServer) ListPub(ctx context.Context, request *articlev1.ListPubRequest) (*articlev1.ListPubResponse, error) {
	arts, err := a.svc.ListPub(ctx, request.GetStartTime().AsTime(), int(request.GetOffset()), int(request.GetLimit()))
	if err != nil {
		return nil, err
	}
	res := make([]*articlev1.Article, 0, len(arts))
	for _, art := range arts {
		res = append(res, a.convertToView(art))
	}
	return &articlev1.ListPubResponse{
		Articles: res,
	}, nil
}

func (a *ArticleServiceServer) convertToDomain(article *articlev1.Article) domain.Article {
	return domain.Article{
		Id:      article.GetId(),
		Title:   article.GetTitle(),
		Content: article.GetContent(),
		Author: domain.Author{
			Id:   article.GetAuthor().GetId(),
			Name: article.GetAuthor().GetName(),
		},
		Status: domain.ArticleStatus(article.GetStatus()),
		Ctime:  article.GetCtime().AsTime(),
		Utime:  article.GetUtime().AsTime(),
	}
}

func (a *ArticleServiceServer) convertToView(art domain.Article) *articlev1.Article {
	return &articlev1.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: &articlev1.Author{
			Id:   art.Author.Id,
			Name: art.Author.Name,
		},
		Status: int32(art.Status),
		Ctime:  timestamppb.New(art.Ctime),
		Utime:  timestamppb.New(art.Utime),
	}
}
