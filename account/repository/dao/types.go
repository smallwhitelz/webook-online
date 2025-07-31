package dao

import "context"

type AccountDAO interface {
	AddActivities(ctx context.Context, activities ...AccountActivity) error
}

// Account 账号本体
// 包含了当前状态
type Account struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 这个账号是哪个用户的账号，如果是系统账号，这个没有取值
	Uid int64
	// 账号ID 这个才是对外使用的
	// 唯一标识一个账号
	Account int64 `gorm:"uniqueIndex:account_type"`
	// 一个人可能有多个账号，可以在这里进行区别
	Type uint8 `gorm:"uniqueIndex:account_type"`

	// 账号本身可以有很多额外的字段
	// 例如和会计有关的，跟税务有关的，跟审计有关的，跟安全有关的

	// 可用余额
	// 一般来说 一种货币就是一个账号，比较好处理
	// 有些一个账号，但是支持多种货币，那么就需要关联另一张表了
	// 记录每一个币种的余额
	Balance  int64
	Currency string

	Ctime int64
	Utime int64
}

// AccountBank AccountEdit 这里也就是说如果你要接入审计，安全，税务之类的，就可以继续加结构体

// AccountActivity 记录流水
type AccountActivity struct {
	Id  int64 `gorm:"primaryKey,autoIncrement"`
	Uid int64

	// 在 biz, biz_id, account 和 account_id 上创建一个联合唯一索引
	// 这样可以确保记账的时候不会重复记账
	Biz   string `gorm:"index:biz_type_id"`
	BizId int64  `gorm:"index:biz_type_id"`

	Account     int64 `gorm:"Index:account_type"`
	AccountType uint8 `gorm:"Index:account_type"`

	Amount   int64
	Currency string

	Ctime int64
	Utime int64
}
