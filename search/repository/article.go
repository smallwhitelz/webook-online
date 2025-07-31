package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"webook/search/domain"
	"webook/search/repository/dao"
)

type articleSearchRepository struct {
	artDao     dao.ArticleSearchDAO
	tagsDao    dao.TagDAO
	likeDao    dao.LikeDAO
	collectDao dao.CollectDAO
}

func NewArticleSearchRepository(artDao dao.ArticleSearchDAO, tagsDao dao.TagDAO,
	likeDao dao.LikeDAO, collectDao dao.CollectDAO) ArticleSearchRepository {
	return &articleSearchRepository{artDao: artDao, tagsDao: tagsDao, likeDao: likeDao, collectDao: collectDao}
}

func (a *articleSearchRepository) SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error) {
	var eg errgroup.Group
	var collectArtIds, artIds, likeArtIds []int64
	var err error
	eg.Go(func() error {
		artIds, err = a.tagsDao.Search(ctx, uid, "article", keywords)
		return err
	})
	eg.Go(func() error {
		likeArtIds, err = a.likeDao.Search(ctx, uid, "article")
		return err
	})
	eg.Go(func() error {
		collectArtIds, err = a.collectDao.Search(ctx, uid, "article")
		return err
	})
	if err = eg.Wait(); err != nil {
		return nil, err
	}
	arts, err := a.artDao.Search(ctx, dao.SearchReq{
		ArtIds:     artIds,
		LikeIds:    likeArtIds,
		CollectIds: collectArtIds,
	}, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(arts, func(idx int, src dao.Article) domain.Article {
		return domain.Article{
			Id:      src.Id,
			Title:   src.Title,
			Content: src.Content,
			Status:  src.Status,
			Tags:    src.Tags,
		}
	}), nil
}

func (a *articleSearchRepository) InputArticle(ctx context.Context, art domain.Article) error {
	return a.artDao.InputArticle(ctx, dao.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  art.Status,
	})
}
