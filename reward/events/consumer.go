package events

import (
	"context"
	"github.com/IBM/sarama"
	"strings"
	"time"
	"webook/pkg/logger"
	"webook/pkg/saramax"
	"webook/reward/domain"
	"webook/reward/service"
)

// PaymentEvent
// 微服务架构reward肯定用不到payment的 PaymentEvent
type PaymentEvent struct {
	BizTradeNO string
	Status     uint8
}

func (p PaymentEvent) ToDomainStatus() domain.RewardStatus {
	// 	PaymentStatusInit
	//	PaymentStatusSuccess
	//	PaymentStatusFailed
	//	PaymentStatusRefund
	switch p.Status {
	case 1:
		return domain.RewardStatusInit
	case 2:
		return domain.RewardStatusPayed
	case 3, 4:
		return domain.RewardStatusFailed
	default:
		return domain.RewardStatusUnknown
	}
}

type PaymentEventConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.RewardService
}

func NewPaymentEventConsumer(client sarama.Client, l logger.LoggerV1, svc service.RewardService) *PaymentEventConsumer {
	return &PaymentEventConsumer{client: client, l: l, svc: svc}
}

func (r *PaymentEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("reward", r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"payment_events"}, saramax.NewHandler[PaymentEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("reward退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *PaymentEventConsumer) Consume(msg *sarama.ConsumerMessage, evt PaymentEvent) error {
	// 判断这个消息的BizTradeNO是否是我们的
	// 因为payment支付那边，这个消费者只是reward大赏的消费者
	// 有可能也会有其他订单类型的消费者去消费payment生产的支付结果消息
	// 所以要进行判断
	if !strings.HasPrefix(evt.BizTradeNO, "reward") {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return r.svc.UpdateReward(ctx, evt.BizTradeNO, evt.ToDomainStatus())
}
