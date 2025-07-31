package events

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/feed/domain"
	"webook/feed/service"
	"webook/pkg/logger"
	"webook/pkg/saramax"
)

const topicFeedEvent = "feed_event"

// FeedEvent 业务方就按照这个格式，将放到feed里面的数据，丢到feed_event这个topic下
type FeedEvent struct {
	Type string
	// 一定是string
	// map[string]any
	// 传过来可能是int64，再反解析回来，就变成了float64
	Metadata map[string]string
}

type FeedEventConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.FeedService
}

func NewFeedEventConsumer(client sarama.Client, l logger.LoggerV1, svc service.FeedService) *FeedEventConsumer {
	return &FeedEventConsumer{client: client, l: l, svc: svc}
}

func (f *FeedEventConsumer) Start() error {
	sg, err := sarama.NewConsumerGroupFromClient("feed_event", f.client)
	if err != nil {
		return err
	}
	go func() {
		err := sg.Consume(context.Background(), []string{topicFeedEvent}, saramax.NewHandler[FeedEvent](f.l, f.Consume))
		if err != nil {
			f.l.Error("feed流退出消费循环", logger.Error(err))
		}
	}()
	return err
}

func (f *FeedEventConsumer) Consume(msg *sarama.ConsumerMessage, event FeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return f.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: event.Type,
		Ext:  event.Metadata,
	})
}
