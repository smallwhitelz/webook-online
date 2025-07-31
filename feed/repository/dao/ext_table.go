package dao

// Feed流中使用扩展表的方式
// 缺点在于代码写起来特别麻烦，工作量是好几倍
type FeedEvent struct {
	Id int64
	// 标注一些类型
	Type string
	// 公共字段，比如说排序字段
}

type ArticleEvent struct {
	Id int64
	// 指向FeedEvent
	Fid int64

	// 文章id
	Aid int64

	// 可以继续冗余其他字段
	AuthorName string
	Title      string
}
