package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type AnyESDAO struct {
	client *elastic.Client
}

func NewAnyESDAO(client *elastic.Client) AnyDAO {
	return &AnyESDAO{client: client}
}

func (a *AnyESDAO) Delete(ctx context.Context, idxName string, docId string, data string) error {
	_, err := a.client.Delete().Index(idxName).Id(docId).Do(ctx)
	return err
}

func (a *AnyESDAO) Input(ctx context.Context, index string, docId string, data string) error {
	_, err := a.client.Index().
		Index(index).Id(docId).BodyString(data).Do(ctx)
	return err
}
