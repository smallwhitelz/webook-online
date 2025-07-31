package web

type CommentVo struct {
	Id          int64       `json:"id"`
	Commentator User        `json:"user"`
	Biz         string      `json:"biz"`
	BizId       int64       `json:"bizId"`
	Content     string      `json:"content"`
	RootID      int64       `json:"rootId,omitempty"`   // 新增根评论ID
	ParentID    int64       `json:"parentId,omitempty"` // 新增父评论ID
	Children    []CommentVo `json:"children,omitempty"`
	Ctime       string      `json:"ctime"`
	Utime       string      `json:"utime"`
}

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AddComment struct {
	Biz      string `json:"biz"`
	BizId    int64  `json:"bizId"`
	Content  string `json:"content"`
	RootID   int64  `json:"rootId,omitempty"`   // 改为直接使用ID
	ParentID int64  `json:"parentId,omitempty"` // 改为直接使用ID
}

type GetCommentListReq struct {
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	MinId int64  `json:"min_id"`
	Limit int64  `json:"limit"`
}

type GetMoreRepliesReq struct {
	Rid   int64 `json:"rid"`
	MaxId int64 `json:"max_id"`
	Limit int64 `json:"limit"`
}
