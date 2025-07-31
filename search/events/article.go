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

const topicSyncArticle = "sync_article_event"

type ArticleEvent struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  int32  `json:"status"`
}

type ArticleConsumer struct {
	artSvc service.SyncService
	client sarama.Client
	l      logger.LoggerV1
}

func NewArticleConsumer(artSvc service.SyncService, client sarama.Client, l logger.LoggerV1) *ArticleConsumer {
	return &ArticleConsumer{artSvc: artSvc, client: client, l: l}
}

func (a *ArticleConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_article", a.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{topicSyncArticle}, saramax.NewHandler[ArticleEvent](a.l, a.Consume))
		if err != nil {
			a.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (a *ArticleConsumer) Consume(msg *sarama.ConsumerMessage, evt ArticleEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return a.artSvc.InputArticle(ctx, domain.Article{
		Id:      evt.Id,
		Title:   evt.Title,
		Content: evt.Content,
		Status:  evt.Status,
	})
}
