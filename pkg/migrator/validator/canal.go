package validator

import (
	"context"
	"gorm.io/gorm"
	"time"
	"webook/pkg/logger"
	"webook/pkg/migrator"
	"webook/pkg/migrator/events"
)

// CanalIncrValidator 使用canal进行增量的数据校验
// 这里用在数据迁移
type CanalIncrValidator[T migrator.Entity] struct {
	// 数据迁移，肯定有源头和目标
	base     *gorm.DB
	target   *gorm.DB
	l        logger.LoggerV1
	producer events.Producer

	direction string
}

func NewCanalIncrValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB,
	l logger.LoggerV1, producer events.Producer, direction string) *CanalIncrValidator[T] {
	return &CanalIncrValidator[T]{base: base, target: target, l: l, producer: producer, direction: direction}
}

// Validate 一次校验一条
// id 是被修改后的数据的主键
func (v *CanalIncrValidator[T]) Validate(ctx context.Context, id int64) error {
	var base T
	err := v.base.WithContext(ctx).Where("id = ?", id).First(&base).Error
	switch err {
	case nil:
		// 找到了
		var target T
		err = v.target.WithContext(ctx).Where("id = ?", id).First(&target).Error
		switch err {
		case nil:
			// 两边都找到了
			if !base.CompareTo(target) {
				v.notify(id, events.InconsistentEventTypeNeq)
			}
			return nil
		case gorm.ErrRecordNotFound:
			// base 有，target没有
			v.notify(id, events.InconsistentEventTypeTargetMissing)
			return nil
		default:
			return err
		}
	case gorm.ErrRecordNotFound:
		// base没找到
		// target也要找找
		var target T
		err = v.target.WithContext(ctx).Where("id = ?", id).First(&target).Error
		switch err {
		case nil:
			// target找到了
			v.notify(id, events.InconsistentEventTypeBaseMissing)
			return nil
		case gorm.ErrRecordNotFound:
			// 两边都没有找到
			return nil
		default:
			return err
		}
	default:
		// 不知道什么错误
		return err
	}
	return nil
}

// notify 通知一条消息
func (v *CanalIncrValidator[T]) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := v.producer.ProduceInconsistentEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Direction: v.direction,
		Type:      typ,
	})
	if err != nil {
		v.l.Error("发送不一致消息失败",
			logger.Int64("id", id),
			logger.String("type", typ),
			logger.Error(err))
	}
}
