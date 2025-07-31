package events

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/follow/repository/cache"
	"webook/follow/repository/dao"
	"webook/pkg/canalx"
	"webook/pkg/logger"
	"webook/pkg/saramax"
)

type MysqlBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	// 直接操作缓存
	// 这里我们要完成的就是：监听 Canal。如果 A 关注了 B，那么就增加 A 的关注人数，并且增加 B 的粉丝数。
	// 维持住 Redis 中的粉丝数量和关注了多少人的数量。
	cache cache.FollowCache
}

func NewMysqlBinlogConsumer(client sarama.Client, l logger.LoggerV1, cache cache.FollowCache) *MysqlBinlogConsumer {
	return &MysqlBinlogConsumer{client: client, l: l, cache: cache}
}

func (m *MysqlBinlogConsumer) Start() error {
	sg, err := sarama.NewConsumerGroupFromClient("follow_canal", m.client)
	if err != nil {
		return err
	}
	go func() {
		err := sg.Consume(context.Background(), []string{"webook_binlog"},
			saramax.NewHandler[canalx.Message[FollowRelation]](m.l, m.Consume))
		if err != nil {
			m.l.Error("follow业务中维护关注人数和粉丝人数的消费者退出循环", logger.Error(err))
		}
	}()
	return err
}

func (m *MysqlBinlogConsumer) Consume(msg *sarama.ConsumerMessage, val canalx.Message[FollowRelation]) error {
	// 先判断是否是我们要处理的数据
	if val.Table != "follow_relations" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, data := range val.Data {
		var err error
		switch data.Status {
		case dao.FollowRelationStatusActive:
			err = m.cache.Follow(ctx, data.Follower, data.Followee)
		case dao.FollowRelationStatusInactive:
			err = m.cache.CancelFollow(ctx, data.Follower, data.Followee)
		default:
			m.l.Error("状态有问题")
		}
		if err != nil {
			m.l.Error("更新关注数和粉丝数失败", logger.Error(err))
		}
	}
	return nil
}

type FollowRelation struct {
	ID int64 `gorm:"column:id;autoIncrement;primaryKey;"`

	Follower int64 `gorm:"uniqueIndex:follower_followee"`
	Followee int64 `gorm:"uniqueIndex:follower_followee"`

	// 软删除策略
	Status uint8

	Ctime int64
	Utime int64
}
