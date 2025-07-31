package repository

import (
	"context"
	"webook/search/repository/dao"
)

type AnyRepository interface {
	Input(ctx context.Context, index string, docId string, data string) error
	Delete(ctx context.Context, idxName string, docId string, data string) error
}

type anyRepository struct {
	dao dao.AnyDAO
}

func NewAnyRepository(dao dao.AnyDAO) AnyRepository {
	return &anyRepository{dao: dao}
}

func (a *anyRepository) Delete(ctx context.Context, idxName string, docId string, data string) error {
	return a.dao.Delete(ctx, idxName, docId, data)
}

func (a *anyRepository) Input(ctx context.Context, index string, docId string, data string) error {
	return a.dao.Input(ctx, index, docId, data)
}
