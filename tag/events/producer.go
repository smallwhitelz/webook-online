package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type SyncDataEvent struct {
	IndexName string
	DocID     string
	// 这里应该是 BizTags
	Data string
}

type Producer interface {
	ProduceSyncEvent(ctx context.Context, data BizTags) error
}

type SaramaSyncProducer struct {
	client sarama.SyncProducer
}

func NewSaramaSyncProducer(client sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{client: client}
}

func (s *SaramaSyncProducer) ProduceSyncEvent(ctx context.Context, data BizTags) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	evt := SyncDataEvent{
		IndexName: "tags_index",
		DocID:     fmt.Sprintf("%d_%s_%d", data.Uid, data.Biz, data.BizId),
		Data:      string(val),
	}
	val, err = json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.client.SendMessage(&sarama.ProducerMessage{
		Topic: "sync_search_data",
		Value: sarama.ByteEncoder(val),
	})
	return err
}

type BizTags struct {
	Tags  []string `json:"tags"`
	Biz   string   `json:"biz"`
	BizId int64    `json:"biz_id"`
	Uid   int64    `json:"uid"`
}
