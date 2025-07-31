package service

import (
	"context"
	"webook/search/domain"
	"webook/search/repository"
)

type SyncService interface {
	InputUser(ctx context.Context, user domain.User) error
	InputArticle(ctx context.Context, art domain.Article) error
	// InputAny 这个接口要的就是不管你传过来什么，我都可以接入
	// index 你的索引名称
	// docId 操作哪个文档
	// data 你的数据
	InputAny(ctx context.Context, index, docId, data string) error
	// Delete any支持取消点赞
	Delete(ctx context.Context, idxName string, docId string, data string) error
}

type syncService struct {
	userRepo repository.UserSearchRepository
	artRepo  repository.ArticleSearchRepository
	anyRepo  repository.AnyRepository
}

func NewSyncService(userRepo repository.UserSearchRepository,
	artRepo repository.ArticleSearchRepository, anyRepo repository.AnyRepository) SyncService {
	return &syncService{userRepo: userRepo, artRepo: artRepo, anyRepo: anyRepo}
}

func (s *syncService) Delete(ctx context.Context, idxName string, docId string, data string) error {
	return s.anyRepo.Delete(ctx, idxName, docId, data)
}

func (s *syncService) InputAny(ctx context.Context, index, docId, data string) error {
	return s.anyRepo.Input(ctx, index, docId, data)
}

func (s *syncService) InputArticle(ctx context.Context, art domain.Article) error {
	return s.artRepo.InputArticle(ctx, art)
}

func (s *syncService) InputUser(ctx context.Context, user domain.User) error {
	return s.userRepo.InputUser(ctx, user)
}
