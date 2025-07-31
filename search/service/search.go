package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"strings"
	"webook/search/domain"
	"webook/search/repository"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type searchService struct {
	userRepo repository.UserSearchRepository
	artRepo  repository.ArticleSearchRepository
}

func NewSearchService(userRepo repository.UserSearchRepository, artRepo repository.ArticleSearchRepository) SearchService {
	return &searchService{userRepo: userRepo, artRepo: artRepo}
}

func (s *searchService) Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error) {
	// 这里你要搜索用户也要搜索art
	// 要对expression进行处理，生成预解析
	// 在大型搜索平台，这里要做的事情非常多
	// 包括对搜索内容的排序、过滤、计算等
	// 这里我们就简单处理一下
	// 输入预处理
	keywords := strings.Split(expression, " ")
	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		users, err := s.userRepo.SearchUser(ctx, keywords)
		res.Users = users
		return err
	})

	eg.Go(func() error {
		arts, err := s.artRepo.SearchArticle(ctx, uid, keywords)
		res.Articles = arts
		return err
	})
	return res, eg.Wait()
}
