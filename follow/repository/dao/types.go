package dao

import "context"

type FollowDAO interface {
	CreateFollowRelation(ctx context.Context, fr FollowRelation) error
	UpdateStatus(ctx context.Context, fr FollowRelation, status uint8) error
	GetFolloweeList(ctx context.Context, follower int64, offset int64, limit int64) ([]FollowRelation, error)
	FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error)
	CntFollower(ctx context.Context, uid int64) (int64, error)
	CntFollowee(ctx context.Context, uid int64) (int64, error)
	GetFollowerList(ctx context.Context, followee int64, offset int64, limit int64) ([]FollowRelation, error)
}

type FollowRelation struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	// 要在这两个列上，创建一个联合唯一索引
	// 这是因为在底层直接就禁止重复关注
	// 如果你认为查询一个人关注了多少人，是主要查询场景
	// <follower, followee>
	// 如果你认为查询一个人有哪些粉丝，是主要查询场景
	// <followee, follower>
	// 我查我关注了哪些人？ WHERE follower = 123(我的 uid)
	Follower int64 `gorm:"uniqueIndex:follower_followee"`
	Followee int64 `gorm:"uniqueIndex:follower_followee"`

	// 软删除策略
	// 也就是说你取消关注某个人的时候，采用软删除
	Status uint8

	// 如果你的关注有类型的，有优先级，有一些备注数据的
	// Type string
	// Priority string
	// Gid 分组ID

	Ctime int64
	Utime int64
}

const (
	FollowRelationStatusUnknown uint8 = iota
	// 关注
	FollowRelationStatusActive
	// 取消关注
	FollowRelationStatusInactive
)

// FollowStatics 这里不是展示列表
// 主要是展示打开个人主页默认的数字
type FollowStatics struct {
	Id  int64 `gorm:"primaryKey,autoIncrement"`
	Uid int64 `gorm:"unique"`
	// 有多少粉丝
	Followers int64
	// 关注了多少人
	Followees int64

	Ctime int64
	Utime int64
}

// UserRelation UserRelationV1 下面的这两个方案指的是要不要将关注、拉黑、屏蔽放到一张表中
// 不要这么玩，最好就是这三个分别是三张表
// 因为你关注的操作一定比其他两个操作更加频繁
type UserRelation struct {
	ID     int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid1   int64 `gorm:"column:uid1;type:int(11);not null;uniqueIndex:user_contact_index"`
	Uid2   int64 `gorm:"column:uid2;type:int(11);not null;uniqueIndex:user_contact_index"`
	Block  bool  // 拉黑
	Mute   bool  // 屏蔽
	Follow bool  // 关注
}

type UserRelationV1 struct {
	ID   int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid1 int64 `gorm:"column:uid1;type:int(11);not null;uniqueIndex:user_contact_index"`
	Uid2 int64 `gorm:"column:uid2;type:int(11);not null;uniqueIndex:user_contact_index"`
	Type string
}
