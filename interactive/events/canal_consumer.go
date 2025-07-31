package events

import (
	"context"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
	"webook/pkg/canalx"
	"webook/pkg/logger"
	"webook/pkg/migrator"
	"webook/pkg/migrator/events"
	"webook/pkg/migrator/validator"
	"webook/pkg/saramax"
)

// MysqlBinlogConsumer
// 这个是利用canal进行数据迁移中的校验部分
type MysqlBinlogConsumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.LoggerV1
	srcToDst *validator.CanalIncrValidator[T]
	dstToSrc *validator.CanalIncrValidator[T]
	dstFirst *atomic.Bool
}

func NewMysqlBinlogConsumer[T migrator.Entity](l logger.LoggerV1, p events.Producer,
	src *gorm.DB, dst *gorm.DB, client sarama.Client) *MysqlBinlogConsumer[T] {
	srcToDst := validator.NewCanalIncrValidator[T](src, dst, l, p, "SRC")
	dstToSrc := validator.NewCanalIncrValidator[T](dst, src, l, p, "DST")
	return &MysqlBinlogConsumer[T]{
		l:        l,
		client:   client,
		srcToDst: srcToDst,
		dstToSrc: dstToSrc,
	}
}

func (m *MysqlBinlogConsumer[T]) Start() error {
	sg, err := sarama.NewConsumerGroupFromClient("migrator_incr", m.client)
	if err != nil {
		return err
	}
	go func() {
		err := sg.Consume(context.Background(), []string{"webook_binlog"}, saramax.NewHandler[canalx.Message[T]](m.l, m.Consume))
		if err != nil {
			m.l.Error("interactive业务数据校验消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (m *MysqlBinlogConsumer[T]) Consume(msg *sarama.ConsumerMessage, val canalx.Message[T]) error {
	// 首先判断这个消息要不要处理
	// 而且这里我们是有阶段的，而且要看是以源表为准还是目标表为准
	dstFirst := m.dstFirst.Load()
	var v *validator.CanalIncrValidator[T]
	// db:
	//  src:
	//    dsn: "root:root@tcp(localhost:13316)/webook"
	//  dst:
	//    dsn: "root:root@tcp(localhost:13316)/webook_intr"
	// 出于保险，你可以进一步校验表名
	if dstFirst && val.Database == "webook_intr" {
		// 以目标表为准，过来的也恰好是目标表的 binlog
		// 要校验
		v = m.dstToSrc
	} else if !dstFirst && val.Database == "webook" {
		// 以源表为准，过来的也恰好是源表的 binlog
		// 要校验
		v = m.srcToDst
	}
	if v == nil {
		return nil
	}
	for _, data := range val.Data {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := v.Validate(ctx, data.ID())
		cancel()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MysqlBinlogConsumer[T]) DstFirst() {
	m.dstFirst.Store(true)
}
