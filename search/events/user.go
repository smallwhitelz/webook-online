package events

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/pkg/logger"
	"webook/pkg/saramax"
	"webook/search/domain"
	"webook/search/service"
)

const topicSyncUser = "sync_user_event"

type UserEvent struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type UserConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       logger.LoggerV1
}

func NewUserConsumer(syncSvc service.SyncService, client sarama.Client, l logger.LoggerV1) *UserConsumer {
	return &UserConsumer{syncSvc: syncSvc, client: client, l: l}
}

func (u *UserConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_user", u.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{topicSyncUser}, saramax.NewHandler[UserEvent](u.l, u.Consume))
		if err != nil {
			u.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (u *UserConsumer) Consume(msg *sarama.ConsumerMessage, evt UserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return u.syncSvc.InputUser(ctx, domain.User{
		Id:       evt.Id,
		Nickname: evt.Nickname,
		Email:    evt.Email,
		Phone:    evt.Phone,
	})
}
