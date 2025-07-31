package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type CommentDAO interface {
	Insert(ctx context.Context, cm Comment) error
	Delete(ctx context.Context, cm Comment) error
	// FindCommentByBiz 只查找一级评论
	FindCommentByBiz(ctx context.Context, biz string, bizId int64, minID int64, limit int64) ([]Comment, error)
	FindRepliesByPid(ctx context.Context, id int64, offset int, limit int) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, maxId int64, limit int64) ([]Comment, error)
}

type CommentGORMDAO struct {
	db *gorm.DB
}

func NewCommentGORMDAO(db *gorm.DB) CommentDAO {
	return &CommentGORMDAO{db: db}
}

func (c *CommentGORMDAO) FindRepliesByRid(ctx context.Context, rid int64, maxId int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).Where("root_id = ? AND id > ?", rid, maxId).
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error
	return res, err
}

// FindRepliesByPid 查找评论的直接的三条评论
func (c *CommentGORMDAO) FindRepliesByPid(ctx context.Context, id int64, offset int, limit int) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).Where("pid = ? ", id).
		Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (c *CommentGORMDAO) FindCommentByBiz(ctx context.Context, biz string, bizId int64, minID int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND id < ? AND pid IS NULL", biz, bizId, minID).
		Limit(int(limit)).Find(&res).Error
	return res, err
}

func (c *CommentGORMDAO) Insert(ctx context.Context, cm Comment) error {
	return c.db.WithContext(ctx).Create(&cm).Error
}

// Delete 级联删除示例
/**
 假设传入的Id为1，且存在以下记录：
记录A：Id=1, PID=NULL（根评论）
记录B：Id=2, PID=1（子评论，引用记录A）
记录C：Id=3, PID=2（子评论，引用记录B）
执行结果：
删除Id为1的记录时，会级联删除Id为2的记录（因为PID=1）。
删除Id为2的记录时，会级联删除Id为3的记录（因为PID=2）。
最终，记录1、2、3都会被删除。
*/
func (c *CommentGORMDAO) Delete(ctx context.Context, cm Comment) error {
	return c.db.WithContext(ctx).Delete(&Comment{
		Id: cm.Id,
	}).Error
}

// Comment 评论的表结构
type Comment struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 发表评论的人
	// 也就是说，如果你要查询某个人发表的所有评论，就需要在这里建立索引
	Uid int64
	// 评价了什么东西
	Biz   string `gorm:"index:biz_type_id"`
	BizId int64  `gorm:"index:biz_type_id"`
	// 评论的内容
	Content string

	// root_id和pid不一定会有值，所以使用sql.NullInt64
	// 也就是说如果这个字段是NULL，他是根评论
	// 我的根评论是哪一个
	// 这是一个冗余字段，缺了这一个整个comment还是完整的
	RootID sql.NullInt64 `gorm:"column:root_id;index"`
	// 评论的父id，也就是你评论的哪一条评论
	// 也就是说如果这个字段是NULL，他是根评论
	PID sql.NullInt64 `gorm:"column:pid;index"`

	// 外键
	// 外键使用的是PID，关联外键的字段是ID 即其他记录的PID字段指向该记录的ID
	// OnDelete:CASCADE 删除策略采用级联删除
	ParentComment *Comment `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"`
	Ctime         int64
	// 事实上，大部分平台是不允许修改评论的
	Utime int64
}

// TableName 指定创建的时候的表名
func (*Comment) TableName() string {
	return "comments"
}
