package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

const LikeIndexName = "like_index"

type LikeESDAO struct {
	client *elastic.Client
}

func NewLikeESDAO(client *elastic.Client) LikeDAO {
	return &LikeESDAO{client: client}
}

func (l *LikeESDAO) Search(ctx context.Context, uid int64, biz string) ([]int64, error) {
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("uid", uid),
		elastic.NewTermQuery("biz", biz),
	)
	resp, err := l.client.Search(LikeIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]int64, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ie InteractiveEvent
		err := json.Unmarshal(hit.Source, &ie)
		if err != nil {
			return nil, err
		}
		res = append(res, ie.BizId)
	}
	return res, nil
}
