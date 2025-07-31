package dao

import "context"

type UserSearchDAO interface {
	InputUser(ctx context.Context, user User) error
	Search(ctx context.Context, keywords []string) ([]User, error)
}

type ArticleSearchDAO interface {
	InputArticle(ctx context.Context, article Article) error
	Search(ctx context.Context, req SearchReq, keywords []string) ([]Article, error)
}

type AnyDAO interface {
	Input(ctx context.Context, index string, docId string, data string) error
	Delete(ctx context.Context, idxName string, docId string, data string) error
}

// 以下三个借口都是返回bizid
// 也就是说用户在搜索的时候
// 这三个会为搜索返回的默认条件
// 对于个人搜索来说
// 你只有点赞和取消点赞以及收藏，所以这里就不用存储数量
// 取消点赞就直接删除掉这个doc
type TagDAO interface {
	Search(ctx context.Context, uid int64, biz string, keywords []string) ([]int64, error)
}

type LikeDAO interface {
	Search(ctx context.Context, uid int64, biz string) ([]int64, error)
}

type CollectDAO interface {
	Search(ctx context.Context, uid int64, biz string) ([]int64, error)
}

type SearchReq struct {
	ArtIds     []int64
	LikeIds    []int64
	CollectIds []int64
}

// InteractiveEvent 用于解析like和collect
type InteractiveEvent struct {
	Biz   string `json:"biz"`
	BizId int64  `json:"bizId"`
	Uid   int64  `json:"uid"`
}
