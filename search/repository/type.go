package repository

import (
	"context"
	"webook/search/domain"
)

type UserSearchRepository interface {
	InputUser(ctx context.Context, user domain.User) error
	SearchUser(ctx context.Context, keywords []string) ([]domain.User, error)
}

type ArticleSearchRepository interface {
	InputArticle(ctx context.Context, art domain.Article) error
	SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error)
}
