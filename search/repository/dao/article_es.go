package dao

import (
	"context"
	"encoding/json"
	"github.com/ecodeclub/ekit/slice"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const ArticleIndexName = "article_index"

type Article struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Status  int32    `json:"status"`
	Tags    []string `json:"tags"`
}

type ArticleESDAO struct {
	client *elastic.Client
}

func NewArticleESDAO(client *elastic.Client) ArticleSearchDAO {
	return &ArticleESDAO{client: client}
}

func (a *ArticleESDAO) Search(ctx context.Context, req SearchReq, keywords []string) ([]Article, error) {
	queryString := strings.Join(keywords, " ")
	// 2 => published
	status := elastic.NewTermsQuery("status", 2)
	title := elastic.NewMatchQuery("title", queryString).Boost(4)
	content := elastic.NewMatchQuery("content", queryString).Boost(4)
	// 结合tag的写法
	// boost默认1.0 越高意味着相关性越高
	tag := elastic.NewTermsQuery("id", slice.Map(req.ArtIds, func(idx int, src int64) any {
		return src
	})).Boost(2)
	like := elastic.NewTermsQuery("id", slice.Map(req.LikeIds, func(idx int, src int64) any {
		return src
	})).Boost(2)
	// 这里我们做的需求是收藏给予最高的权重
	collect := elastic.NewTermsQuery("id", slice.Map(req.CollectIds, func(idx int, src int64) any {
		return src
	})).Boost(100)
	or := elastic.NewBoolQuery().Should(title, content, tag, like, collect)
	query := elastic.NewBoolQuery().Must(status, or)
	sort := elastic.NewFieldSort("id").Desc()
	scoreSort := elastic.NewFieldSort("_score").Desc()
	resp, err := a.client.Search(ArticleIndexName).
		SortBy(scoreSort, sort).
		Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]Article, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var art Article
		err := json.Unmarshal(hit.Source, &art)
		if err != nil {
			return nil, err
		}
		res = append(res, art)
	}
	return res, nil
}

func (a *ArticleESDAO) InputArticle(ctx context.Context, article Article) error {
	_, err := a.client.Index().Index(ArticleIndexName).
		Id(strconv.FormatInt(article.Id, 10)).
		BodyJson(article).
		Do(ctx)
	return err
}
