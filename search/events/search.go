package events

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/pkg/logger"
	"webook/pkg/saramax"
	"webook/search/service"
)

// SyncDataEvent 通用的 sync data event
// 所有的业务方都可以通过这个event同步数据
type SyncDataEvent struct {
	IndexName string
	DocId     string
	Data      string
}

type SyncDataEventConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       logger.LoggerV1
}

func NewSyncDataEventConsumer(syncSvc service.SyncService, client sarama.Client, l logger.LoggerV1) *SyncDataEventConsumer {
	return &SyncDataEventConsumer{syncSvc: syncSvc, client: client, l: l}
}

func (s *SyncDataEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_data", s.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"sync_search_data"}, saramax.NewHandler[SyncDataEvent](s.l, s.Consume))
		if err != nil {
			s.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (s *SyncDataEventConsumer) Consume(msg *sarama.ConsumerMessage, evt SyncDataEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return s.syncSvc.InputAny(ctx, evt.IndexName, evt.DocId, evt.Data)
}
