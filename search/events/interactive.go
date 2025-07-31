package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"webook/pkg/logger"
	"webook/pkg/saramax"
	"webook/search/service"
)

const interactiveTopic = "interactive_sync"

type InteractiveEvent struct {
	Type  int64  `json:"type,omitempty"` // 1-喜欢 2-收藏 3-取消喜欢
	Biz   string `json:"biz"`
	BizId int64  `json:"bizId"`
	Uid   int64  `json:"uid"`
}

type InteractiveConsumer struct {
	client  sarama.Client
	syncSvc service.SyncService
	l       logger.LoggerV1
}

func NewInteractiveConsumer(client sarama.Client, syncSvc service.SyncService, l logger.LoggerV1) *InteractiveConsumer {
	return &InteractiveConsumer{client: client, syncSvc: syncSvc, l: l}
}

func (i *InteractiveConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_interactive_group", i.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{interactiveTopic}, saramax.NewHandler[InteractiveEvent](i.l, i.Consume))
		if err != nil {
			i.l.Error("点赞收藏接入搜索退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (i *InteractiveConsumer) Consume(msg *sarama.ConsumerMessage, evt InteractiveEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	switch evt.Type {
	case 1:
		return i.syncSvc.InputAny(context.Background(), "like_index", i.getDocId(evt), string(val))
	case 2:
		return i.syncSvc.InputAny(context.Background(), "collect_index", i.getDocId(evt), string(val))
	case 3:
		return i.syncSvc.Delete(context.Background(), "like_index", i.getDocId(evt), string(val))
	default:
		return errors.New("未知参数")
	}
}

func (i *InteractiveConsumer) getDocId(data InteractiveEvent) string {
	return fmt.Sprintf("%d_%s_%d", data.Uid, data.Biz, data.BizId)
}
